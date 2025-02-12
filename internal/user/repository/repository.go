package repository

import (
	"Coolshop/internal/connection"
	"Coolshop/internal/model/user"
	def "Coolshop/internal/user"
	"Coolshop/pkg/errlist"
	"context"
	"fmt"
	"time"
)

var _ def.Repository = (*repository)(nil)

type repository struct {
	psqlDB connection.DB
}

func NewRepository(psqlDb connection.DB) *repository {
	return &repository{psqlDB: psqlDb}
}

func (r *repository) Create(ctx context.Context, data user.UserData) error {
	_, err := r.psqlDB.Execute(ctx,
		"INSERT INTO users(name, surname, username, password, email) VALUES ($1,$2,$3, $4, $5);",
		data.Name,
		data.Surname,
		data.Username,
		data.Password,
		data.Email,
	)
	if err != nil {
		return fmt.Errorf("r.psqlDB.ExecContext: %w", err)
	}

	return nil
}

func (r *repository) Find(ctx context.Context, email string) error {
	var exists bool

	err := r.psqlDB.QueryRow(ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)",
		email,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("r.psqlDB.QueryRowContext: %w", err)
	}

	if exists {
		return errlist.ErrAlreadyExists
	}

	return nil
}

func (r *repository) Get(ctx context.Context, email string) (user.UserData, error) {
	var data user.UserData

	err := r.psqlDB.Get(ctx, &data, "SELECT id, name, surname, username, password, email, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL", email)
	if err != nil {
		return user.UserData{}, err
	}

	return data, nil
}

func (r *repository) Delete(ctx context.Context, userID int64) error {
	result, err := r.psqlDB.Execute(
		ctx,
		"UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL",
		time.Now(),
		userID,
	)
	if err != nil {
		return fmt.Errorf("r.psqlDB.ExecContext: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlist.ErrNotFound
	}

	return nil
}

func (r *repository) GetByID(ctx context.Context, userID int64) (user.UserData, error) {
	var userData user.UserData

	err := r.psqlDB.Get(
		ctx,
		&userData,
		"SELECT id, name, surname, username, password, email, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL",
		userID,
	)
	if err != nil {
		return user.UserData{}, err
	}

	return userData, nil
}

func (r *repository) UpdatePassword(ctx context.Context, userID int64, newHashedPassword string) error {
	result, err := r.psqlDB.Execute(
		ctx,
		"Update users SET password = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL",
		newHashedPassword,
		userID,
	)
	if err != nil {
		return fmt.Errorf("r.psqlDB.ExecContext: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if rows == 0 {
		return errlist.ErrNotFound
	}

	return nil
}
