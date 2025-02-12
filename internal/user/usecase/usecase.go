package usecase

import (
	"Coolshop/internal/auth"
	"Coolshop/internal/cache"
	userModel "Coolshop/internal/model/user"
	"Coolshop/internal/user"
	"Coolshop/pkg/errlist"
	"Coolshop/pkg/utilities"
	"context"
	"database/sql"
	"errors"
	"strconv"
)

var _ user.UseCase = (*useCase)(nil)

type useCase struct {
	repository user.Repository
	cache      cache.Repository
	jwtGen     auth.JwtGen
}

func NewUseCase(repository user.Repository, cache cache.Repository, jwtGen auth.JwtGen) *useCase {
	return &useCase{
		repository: repository,
		cache:      cache,
		jwtGen:     jwtGen,
	}
}

func (u *useCase) SignUp(ctx context.Context, user userModel.User) error {
	err := u.repository.Find(ctx, user.Email)
	if err != nil {
		return err
	}

	user.Password, err = utilities.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = u.repository.Create(ctx, user.ToStorage())
	if err != nil {
		return err
	}

	return nil
}

func (u *useCase) Login(ctx context.Context, user userModel.UserLogin) (userModel.Tokens, error) {
	userData, err := u.repository.Get(ctx, user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userModel.Tokens{}, errlist.ErrNotFound
		}

		return userModel.Tokens{}, err
	}

	passwordIsValid := utilities.VerifyPassword(userData.Password, user.Password)
	if !passwordIsValid {
		return userModel.Tokens{}, errlist.ErrInvalidPassword
	}

	accessToken, refreshToken, err := u.jwtGen.GenerateJWT(strconv.FormatInt(userData.ID, 10))
	if err != nil {
		return userModel.Tokens{}, err
	}

	err = u.cache.SaveRefreshToken(ctx, strconv.FormatInt(userData.ID, 10), refreshToken)
	if err != nil {
		return userModel.Tokens{}, err
	}

	return userModel.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *useCase) GetNewAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := u.jwtGen.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	storedToken, err := u.cache.GetRefreshToken(ctx, claims.UserID)
	if err != nil || storedToken != refreshToken {
		return "", errors.New("invalid refresh token")
	}

	newAccessToken, _, err := u.jwtGen.GenerateJWT(claims.UserID)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (u *useCase) Logout(ctx context.Context, userIDStr string) error {
	_, err := u.cache.GetRefreshToken(ctx, userIDStr)
	if err != nil {
		if errors.Is(err, errlist.ErrNotFound) {
			// If no token exists, user is already logged out
			return nil
		}

		return err
	}

	err = u.cache.DeleteRefreshToken(ctx, userIDStr)
	if err != nil {
		return err
	}

	return nil
}

func (u *useCase) Delete(ctx context.Context, userIDStr string) error {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return err
	}

	err = u.cache.DeleteRefreshToken(ctx, userIDStr)
	if err != nil {
		return err
	}

	err = u.repository.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *useCase) GetUserByID(ctx context.Context, userIDStr string) (userModel.UserResponse, error) {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return userModel.UserResponse{}, err
	}

	userData, err := u.repository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userModel.UserResponse{}, errlist.ErrNotFound
		}

		return userModel.UserResponse{}, err
	}

	return userData.ToServer(), nil
}

func (u *useCase) UpdatePassword(ctx context.Context, userIDStr string, newPassword string) error {
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return err
	}

	newHashedPassword, err := utilities.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = u.repository.UpdatePassword(ctx, userID, newHashedPassword)
	if err != nil {
		return err
	}

	return nil
}
