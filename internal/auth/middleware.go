package auth

import (
	"Coolshop/internal/model/middleware"
	"Coolshop/logger"
	"Coolshop/pkg/errlist"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtGen JwtGen) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()

			return
		}

		claims, err := jwtGen.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()

			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// CustomErrorHandler is a middleware that formats error responses.
func CustomErrorHandler(log logger.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process request

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			statusCode := getStatusCode(err.Err)

			log.Error("API Error",
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", statusCode),
				zap.String("error", err.Error()),
			)

			c.AbortWithStatusJSON(statusCode, middleware.ErrorResponse{
				Error:   http.StatusText(statusCode),
				Message: err.Err.Error(),
			})
		}
	}
}

// getStatusCode maps error messages from `errlist` to HTTP status codes.
func getStatusCode(err error) int {
	switch {
	case errors.Is(err, errlist.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, errlist.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, errlist.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errlist.ErrInvalidPassword):
		return http.StatusUnauthorized
	case errors.Is(err, errlist.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, errlist.ErrAlreadyExists):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
