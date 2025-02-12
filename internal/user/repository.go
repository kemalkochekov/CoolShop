package user

import (
	"Coolshop/internal/model/user"
	"context"
)

type Repository interface {
	Create(ctx context.Context, data user.UserData) error
	Find(ctx context.Context, email string) error
	Get(ctx context.Context, email string) (user.UserData, error)
	Delete(ctx context.Context, userID int64) error
	GetByID(ctx context.Context, userID int64) (user.UserData, error)
	UpdatePassword(ctx context.Context, userID int64, newHashedPassword string) error
}
