package delivery

import (
	userModel "Coolshop/internal/model/user"
	"Coolshop/internal/user"
	"Coolshop/logger"
	"Coolshop/pkg/constant"
	"Coolshop/pkg/errlist"
	"Coolshop/pkg/reqvalidator"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

var _ user.Handlers = (*handler)(nil)

type handler struct {
	useCase user.UseCase
	logger  logger.LoggerInterface
}

func NewHandler(useCase user.UseCase, logger logger.LoggerInterface) *handler {
	return &handler{useCase: useCase, logger: logger}
}

// HealthCheck returns the status of the service.
//
//	@Summary Health Check
//	@Description Checks if the API is running.
//	@Tags system
//	@Produce json
//	@Success 200 {object} user.StandardResponse "status: healthy"
//	@Router /health [get]
func (h *handler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	}
}

// SignUp creates a new user account.
//
//	@Summary User Registration
//	@Description Register a new user with email and password.
//	@Tags auth
//	@Accept json
//	@Produce json
//	@Param request body user.User true "User registration details"
//	@Success 200 {object} user.StandardResponse "status: ok"
//	@Failure 400 {object} middleware.ErrorResponse "Invalid request"
//	@Failure 409 {object} middleware.ErrorResponse "User already exists"
//	@Failure 500 {object} middleware.ErrorResponse "Internal server error"
//	@Router /auth/sign_up [post]
func (h *handler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := userModel.User{}
		if err := c.BindJSON(&request); err != nil {
			h.logger.Warn("SignUp: Invalid request", zap.Error(err))
			c.Error(errlist.ErrBadRequest)

			return
		}

		err := reqvalidator.ValidateRequest(request)
		if err != nil {
			h.logger.Warn("SignUp: Validation failed", zap.Error(err))
			c.Error(errlist.ErrBadRequest)

			return
		}

		err = h.useCase.SignUp(c, request)
		if err != nil {
			if errors.Is(err, errlist.ErrAlreadyExists) {
				h.logger.Warn("SignUp: User already exists", zap.String("email", request.Email))
				c.Error(errlist.ErrAlreadyExists)

				return
			}

			h.logger.Error("SignUp: Internal error", zap.Error(err))
			c.Error(errlist.ErrInternalServer)

			return
		}

		h.logger.Info("SignUp successful", zap.String("email", request.Email))
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// Login authenticates a user and returns access and refresh tokens.
//
//	@Summary User Login
//	@Description Authenticate user with email and password.
//	@Tags auth
//	@Accept json
//	@Produce json
//	@Param request body user.UserLogin true "User credentials"
//	@Success 200 {object} user.Tokens "Access token is returned in response. Refresh token is set in a cookie."
//	@Failure 400 {object} middleware.ErrorResponse "Bad request"
//	@Failure 404 {object} middleware.ErrorResponse "User not found"
//	@Failure 500 {object} middleware.ErrorResponse "Internal server error"
//	@Router /auth/login [post]
func (h *handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := userModel.UserLogin{}
		if err := c.BindJSON(&request); err != nil {
			h.logger.Warn("Login: Invalid request", zap.Error(err))
			c.Error(errlist.ErrBadRequest)

			return
		}

		err := reqvalidator.ValidateRequest(request)
		if err != nil {
			h.logger.Warn("Login: Validation failed", zap.Error(err))
			c.Error(errlist.ErrBadRequest)

			return
		}

		token, err := h.useCase.Login(c, request)
		if err != nil {
			if errors.Is(err, errlist.ErrNotFound) {
				h.logger.Warn("Login: User not found", zap.String("email", request.Email))
				c.Error(errlist.ErrNotFound)

				return
			}

			if errors.Is(err, errlist.ErrInvalidPassword) {
				h.logger.Warn("Login: Invalid password", zap.String("email", request.Email))
				c.Error(errlist.ErrInvalidPassword)

				return
			}

			h.logger.Error("Login: Internal error", zap.Error(err))
			c.Error(errlist.ErrInternalServer)

			return
		}

		c.SetCookie(
			"refresh_token",
			token.RefreshToken,
			constant.MaxAgeCookie,
			"/auth/refresh",
			"",
			false,
			true,
		)

		h.logger.Info("Login successful", zap.String("email", request.Email))
		c.JSON(http.StatusOK, gin.H{"access_token": token.AccessToken})
	}
}

// RefreshToken generates a new access token using a refresh token.
//
//	@Summary Refresh Access Token
//	@Description Generates a new access token using the provided refresh token.
//	@Tags auth
//	@Accept json
//	@Produce json
//	@Success 200 {object} user.Tokens "New access token is returned in response."
//	@Failure 401 {object} middleware.ErrorResponse "Refresh token required"
//	@Failure 500 {object} middleware.ErrorResponse "Internal server error"
//	@Router /auth/refresh [post]
func (h *handler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			h.logger.Warn("RefreshToken: Missing refresh token")
			c.Error(errlist.ErrUnauthorized)

			return
		}

		newAccessToken, err := h.useCase.GetNewAccessToken(c, refreshToken)
		if err != nil {
			h.logger.Warn("RefreshToken: Invalid refresh token")
			c.Error(errlist.ErrUnauthorized)

			return
		}

		h.logger.Info("RefreshToken: Access token refreshed")
		c.JSON(http.StatusOK, gin.H{"access_token": newAccessToken})
	}
}

// Logout logs out the user by invalidating their session.
//
//	@Summary User Logout
//	@Description Logs out the authenticated user and invalidates their refresh token. Requires an access token in the `Authorization` header.
//	@Tags auth
//	@Security Bearer
//	@Success 200 {object} user.StandardResponse "message: Logged out successfully"
//	@Failure 401 {object} middleware.ErrorResponse "Unauthorized"
//	@Failure 500 {object} middleware.ErrorResponse "Internal server error"
//	@Router /user/logout [post]
func (h *handler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("userID")
		if !exists {
			h.logger.Warn("Logout: Unauthorized access")
			c.Error(errlist.ErrUnauthorized)

			return
		}

		userID, ok := userIDRaw.(string)
		if !ok {
			h.logger.Warn("Logout: Invalid user ID format")
			c.Error(errlist.ErrBadRequest)

			return
		}

		err := h.useCase.Logout(c, userID)
		if err != nil {
			h.logger.Error("Logout: Internal error", zap.Error(err))
			c.Error(errlist.ErrInternalServer)

			return
		}

		h.logger.Info("Logout successful", zap.String("userID", userID))
		c.SetCookie("refresh_token", "", -1, "/", "", true, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// GetUser retrieves user information by ID.
//
//	@Summary Get User
//	@Description Fetches user details using their ID. Requires an access token in the `Authorization` header.
//	@Tags user
//	@Accept json
//	@Produce json
//	@Param id path int true "User ID"
//	@Success 200 {object} user.UserResponse
//	@Failure 400 {object} middleware.ErrorResponse "Bad request"
//	@Failure 401 {object} middleware.ErrorResponse "Unauthorized"
//	@Failure 404 {object} middleware.ErrorResponse "User not found"
//	@Failure 500 {object} middleware.ErrorResponse "Internal server error"
//	@Router /user/{id} [get]
//	@Security Bearer
func (h *handler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			h.logger.Warn("Bad request: User ID is missing")
			c.Error(errlist.ErrUnauthorized)

			return
		}

		userID := c.Param("id")
		if userID == "" {
			h.logger.Warn("User not found", zap.String("userID", userID))
			c.Error(errlist.ErrBadRequest)

			return
		}

		response, err := h.useCase.GetUserByID(c, userID)
		if err != nil {
			if errors.Is(err, errlist.ErrNotFound) {
				h.logger.Warn("User not found", zap.String("userID", userID))
				c.Error(errlist.ErrNotFound)

				return
			}

			h.logger.Error("Failed to get user", zap.Error(err))
			c.Error(errlist.ErrInternalServer)

			return
		}

		h.logger.Info("GetUser successful", zap.String("userID", userID))
		c.JSON(http.StatusOK, response)
	}
}

// DeleteUser removes the authenticated user from the system.
//
//	@Summary Delete User
//	@Description Deletes the currently logged-in user. Requires an access token in the `Authorization` header.
//	@Tags user
//	@Security Bearer
//	@Success 200 {object} user.StandardResponse "message: User deleted successfully"
//	@Failure 401 {object} middleware.ErrorResponse "Unauthorized"
//	@Failure 404 {object} middleware.ErrorResponse "User not found"
//	@Failure 500 {object} middleware.ErrorResponse "Internal server error"
//	@Router /user [delete]
func (h *handler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("userID")
		if !exists {
			h.logger.Warn("DeleteUser: Unauthorized access")
			c.Error(errlist.ErrUnauthorized)

			return
		}

		userID, ok := userIDRaw.(string)
		if !ok {
			h.logger.Warn("DeleteUser: Invalid user ID format")
			c.Error(errlist.ErrBadRequest)

			return
		}

		err := h.useCase.Delete(c, userID) // Now using `userID` from context
		if err != nil {
			if errors.Is(err, errlist.ErrNotFound) {
				h.logger.Warn("DeleteUser: User not found", zap.String("userID", userID))
				c.Error(errlist.ErrNotFound)

				return
			}

			h.logger.Error("DeleteUser: Internal error", zap.Error(err))
			c.Error(errlist.ErrInternalServer)

			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

// UpdatePassword updates the password for the authenticated user.
//
//	@Summary Update Password
//	@Description Allows an authenticated user to update their password. Requires an access token in the `Authorization` header.
//	@Tags user
//	@Security Bearer
//	@Accept json
//	@Produce json
//	@Param request body user.UpdateRequest true "New password details"
//	@Success 200 {object} user.StandardResponse "message: Password updated successfully"
//	@Failure 400 {object} middleware.ErrorResponse "Invalid request"
//	@Failure 401 {object} middleware.ErrorResponse "Unauthorized"
//	@Failure 500 {object} middleware.ErrorResponse "Internal server error"
//	@Router /user/update_password [put]
func (h *handler) UpdatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get("userID")
		if !exists {
			h.logger.Warn("UpdatePassword: Unauthorized access")
			c.Error(errlist.ErrUnauthorized)

			return
		}

		userIDStr, ok := userIDRaw.(string)
		if !ok {
			h.logger.Warn("UpdatePassword: Invalid user ID format")
			c.Error(errlist.ErrBadRequest)

			return
		}

		updateRequest := userModel.UpdateRequest{}
		if err := c.BindJSON(&updateRequest); err != nil {
			h.logger.Warn("UpdatePassword: Invalid request format", zap.Error(err))
			c.Error(errlist.ErrBadRequest)

			return
		}

		err := reqvalidator.ValidateRequest(updateRequest)
		if err != nil {
			h.logger.Warn("UpdatePassword: Validation failed", zap.Error(err))
			c.Error(errlist.ErrBadRequest)

			return
		}

		err = h.useCase.UpdatePassword(c, userIDStr, updateRequest.Password)
		if err != nil {
			h.logger.Error("UpdatePassword: Internal error", zap.Error(err))
			c.Error(errlist.ErrInternalServer)

			return
		}

		h.logger.Info("UpdatePassword successful", zap.String("userID", userIDStr))
		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}
}
