package router

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(auth_router *gin.RouterGroup) {
	auth_router.POST("/signin", service.UserLogin)
	auth_router.POST("/signup", service.UserRegister)
}

func AddProjectRouter(project_router *gin.RouterGroup) {
	project_router.GET("/list", service.ListAllProjects)
	project_router.POST("/create", service.CreateProject)
	project_router.POST("/:id/addRepo", service.AddRepo)
	project_router.POST("/:id/delete", service.AddRepo)
	project_router.POST("/:id/deleteRepo", service.AddRepo)
}
