package user

import "time"

// User represents a user in the system.
//
//	@Description User model containing user information
type User struct {
	Username string `json:"username" validate:"required" example:"johndoe"`
	Name     string `json:"name" validate:"required" example:"John"`
	Surname  string `json:"surname" validate:"required" example:"Doe"`
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"SuperStrongPassword123"`
}

type UserResponse struct {
	ID        int64     `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Name      string    `json:"name" example:"John"`
	Surname   string    `json:"surname" example:"Doe"`
	Email     string    `json:"email" example:"user@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2020-01-01 12:00:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-01-01 12:00:00"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Tokens represents the authentication tokens returned after login.
//
//	@Description Authentication tokens for user sessions.
type Tokens struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"-" swaggerignore:"true"`
}

type UpdateRequest struct {
	Password string `json:"password" validate:"required"`
}

// StandardResponse represents a common response structure.
//
//	@Description Standard API response
type StandardResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

type UserData struct {
	ID        int64     `db:"id"`
	Username  string    `db:"username"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *User) ToStorage() UserData {
	return UserData{
		Username: u.Username,
		Name:     u.Name,
		Surname:  u.Surname,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (u *UserData) ToServer() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Surname:   u.Surname,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
