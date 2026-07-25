package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Zelviarani16/taskflow-api/cmd"
	"github.com/Zelviarani16/taskflow-api/config/database"
	"github.com/Zelviarani16/taskflow-api/handler"
	"github.com/Zelviarani16/taskflow-api/middleware"
	"github.com/Zelviarani16/taskflow-api/repository"
	"github.com/Zelviarani16/taskflow-api/routes"
	"github.com/Zelviarani16/taskflow-api/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// Jika dijalankan dengan flag (--migrate/--seed/--rollback), jalankan
	// perintah tersebut dan keluar tanpa menjalankan server HTTP.
	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	// Perangkaian dependency: repository -> service -> handler.
	// Setara manual dengan DI container / service bindings di Laravel.
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	jwtService := service.NewJWTService()
	authService := service.NewAuthService(userRepo, jwtService)
	projectService := service.NewProjectService(projectRepo)
	taskService := service.NewTaskService(taskRepo, projectRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	taskHandler := handler.NewTaskHandler(taskService)

	authMiddleware := middleware.Authentication(jwtService)

	server := gin.Default()

	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "TaskFlow API is running"})
	})

	routes.Auth(server, authHandler)
	routes.Project(server, authMiddleware, projectHandler, taskHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("server jalan di port %s", port)
	if err := server.Run(":" + port); err != nil {
		log.Fatalf("gagal jalankan server: %v", err)
	}
}
