package routes

import (
	"github.com/Zelviarani16/taskflow-api/handler"
	"github.com/gin-gonic/gin"
)

// Auth mendaftarkan endpoint autentikasi yang bersifat publik (tidak perlu token).
func Auth(route *gin.Engine, authHandler handler.IAuthHandler) {
	authRoutes := route.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}
}
