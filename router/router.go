package router

import (
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(auth_router *gin.RouterGroup) {
	auth_router.POST("/signin", service.UserLogin)
	auth_router.POST("/signup", service.UserRegister)
	auth_router.POST("/forgot", service.UserForgotPassword)
}

func AddProjectRouter(project_router *gin.RouterGroup) {
	project_router.GET("/repos", service.ListAllRepos)
	project_router.POST("/repos", service.AddRepo)
	project_router.POST("/", service.CreateProject)
	project_router.DELETE("/", service.DeleteProject)
}
