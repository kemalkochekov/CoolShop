package user

import (
	"github.com/gin-gonic/gin"
)

type Handlers interface {
	SignUp() gin.HandlerFunc
	Login() gin.HandlerFunc
	RefreshToken() gin.HandlerFunc
	Logout() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	HealthCheck() gin.HandlerFunc
	UpdatePassword() gin.HandlerFunc
}
