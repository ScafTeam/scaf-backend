package router

import (
	"backend/middleware"
	"backend/service"
	"github.com/gin-gonic/gin"
)

func AddAuthRouter(auth_router *gin.RouterGroup) {
	auth_router.POST("/signup/", service.UserRegister)
	auth_router.POST("/forgot/", service.UserForgotPassword)
	auth_router.POST("/signin/", middleware.AuthMiddleware.LoginHandler)
	auth_router.POST("/signout/", middleware.AuthMiddleware.LogoutHandler)
	auth_router.POST("/refresh/", middleware.AuthMiddleware.RefreshHandler)
	auth_router.Use(middleware.AuthMiddleware.MiddlewareFunc())
	{
		auth_router.GET("/hello", middleware.HelloHandler)
	}
}

func AddProjectRouter(project_router *gin.RouterGroup) {
	project_router.GET("/", service.ListAllProjects)
	project_router.GET("/:project_name/", service.GetProject)
	project_router.Use(middleware.MemberCheck())
	{
		project_router.PUT("/:project_name/", service.UpdateProject)
		project_router.POST("/:project_name/member/", service.AddMember)
	}
	project_router.Use(middleware.AuthCheck())
	{
		project_router.POST("/", service.CreateProject)
		project_router.DELETE("/:project_name/", service.DeleteProject)
	}
}

func AddRepoRouter(repo_router *gin.RouterGroup) {
	repo_router.GET("/", service.ListAllRepos)
	repo_router.Use(middleware.MemberCheck())
	{
		repo_router.POST("/", service.AddRepo)
	}
}

func AddKanbanRouter(kanban_router *gin.RouterGroup) {
	kanban_router.GET("/", service.ListKanban)
	kanban_router.Use(middleware.MemberCheck())
	{
		kanban_router.PUT("/", service.AddWorkFlow)
		kanban_router.POST("/", service.AddTask)
		kanban_router.DELETE("/", service.DeleteTask)
	}
}
