package user

import (
	"Coolshop/internal/model/user"
	"context"
)

type UseCase interface {
	SignUp(ctx context.Context, user user.User) error
	Login(ctx context.Context, user user.UserLogin) (user.Tokens, error)
	GetNewAccessToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, userIDStr string) error
	Delete(ctx context.Context, userIDStr string) error
	GetUserByID(ctx context.Context, userIDStr string) (user.UserResponse, error)
	UpdatePassword(ctx context.Context, userIDStr string, newPassword string) error
}
