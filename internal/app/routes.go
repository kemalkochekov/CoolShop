package app

import (
	"Coolshop/internal/auth"
	"Coolshop/internal/user"

	_ "Coolshop/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *App) initializeRoutes(h user.Handlers) {
	a.server.Use(auth.CustomErrorHandler(a.zapLogger))

	authorization := a.server.Group("/auth")
	{
		authorization.POST("/refresh", h.RefreshToken())
		authorization.POST("/sign_up", h.SignUp())
		authorization.POST("/login", h.Login())
	}

	protected := a.server.Group("/user").Use(auth.AuthMiddleware(a.jwtGen))
	{
		protected.POST("/logout", h.Logout())
		protected.GET("/:id", h.GetUser())
		protected.DELETE("/", h.DeleteUser())
		protected.PUT("/update_password", h.UpdatePassword())
	}

	a.server.GET("/health", h.HealthCheck())
	a.server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
