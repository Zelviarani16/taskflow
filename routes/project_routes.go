package routes

import (
	"github.com/Zelviarani16/taskflow-api/handler"
	"github.com/gin-gonic/gin"
)

// Project mendaftarkan endpoint project dan, bersarang di bawahnya, endpoint
// task yang terikat pada satu project (/projects/:id/tasks/...).
// authMiddleware diterapkan ke seluruh grup karena setiap aksi di sini
// membutuhkan user yang sudah login.
func Project(route *gin.Engine, authMiddleware gin.HandlerFunc, projectHandler handler.IProjectHandler, taskHandler handler.ITaskHandler) {
	projectRoutes := route.Group("/api/v1/projects", authMiddleware)
	{
		projectRoutes.POST("", projectHandler.Create)
		projectRoutes.GET("", projectHandler.List)
		projectRoutes.GET("/:id", projectHandler.GetByID)
		projectRoutes.PUT("/:id", projectHandler.Update)
		projectRoutes.DELETE("/:id", projectHandler.Delete)

		taskRoutes := projectRoutes.Group("/:id/tasks")
		{
			taskRoutes.POST("", taskHandler.Create)
			taskRoutes.GET("", taskHandler.List)
			taskRoutes.PUT("/:taskId", taskHandler.Update)
			taskRoutes.DELETE("/:taskId", taskHandler.Delete)
		}
	}
}
