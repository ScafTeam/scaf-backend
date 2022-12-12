package router

import (
	"backend/middleware"
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(auth_router *gin.RouterGroup) {
	auth_router.POST("/signup", service.UserRegister)
	auth_router.POST("/forgot", service.UserForgotPassword)
	auth_router.POST("/signin", middleware.AuthMiddleware.LoginHandler)
	auth_router.POST("/signout", middleware.AuthMiddleware.LogoutHandler)
	auth_router.POST("/refresh", middleware.AuthMiddleware.RefreshHandler)
	auth_router.Use(middleware.AuthMiddleware.MiddlewareFunc())
	{
		auth_router.GET("/hello", middleware.HelloHandler)
	}
}

func AddProjectRouter(project_router *gin.RouterGroup) {
	project_router.Use(middleware.AuthCheck())
	{
		project_router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello World"})
		})
		project_router.POST("/", service.CreateProject)
		project_router.DELETE("/", service.DeleteProject)
	}
	project_router.GET("/repos", service.ListAllRepos)
	project_router.POST("/repos", service.AddRepo)
	project_router.POST("/", service.CreateProject)
	project_router.DELETE("/", service.DeleteProject)
}
