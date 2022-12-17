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
	project_router.GET("/", service.ListAllProjects)
	project_router.Use(middleware.AuthCheck())
	{
		project_router.DELETE("/:project_id", service.DeleteProject)
		project_router.POST("/", service.CreateProject)
	}
}

func AddRepoRouter(repo_router *gin.RouterGroup) {
	repo_router.GET("/", service.ListAllRepos)
	repo_router.Use(middleware.MemberCheck())
	{
		repo_router.POST("/", service.AddRepo)
	}
}
