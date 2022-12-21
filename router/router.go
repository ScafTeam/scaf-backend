package router

import (
	"backend/middleware"
	"backend/service"
	"github.com/gin-gonic/gin"
)

func Init(server *gin.Engine) {
	auth_router := server.Group("/")
	{
		auth_router.POST("/signup/", service.UserRegister)
		auth_router.POST("/forgot/", service.UserForgotPassword)
		auth_router.POST("/signin/", middleware.AuthMiddleware.LoginHandler)
		auth_router.POST("/signout/", middleware.AuthMiddleware.LogoutHandler)
		auth_router.POST("/refresh/", middleware.AuthMiddleware.RefreshHandler)
	}

	project_router := server.Group("/user/:user_email/project")
	{
		project_router.GET("/", service.ListAllProjects)
		project_router.GET("/:project_name/", service.GetProject)

		project_member_router := server.Group("/user/:user_email/project/").Use(middleware.MemberCheck())
		{
			project_member_router.PUT("/:project_name/", service.UpdateProject)
			project_member_router.POST("/:project_name/member/", service.AddMember)
		}

		project_owner_router := server.Group("/user/:user_email/project/").Use(middleware.OwnerCheck())
		{
			project_owner_router.POST("/", service.CreateProject)
			project_owner_router.DELETE("/:project_name/", service.DeleteProject)
		}
	}

	repo_router := project_router.Group("/:project_name/repo")
	{
		repo_router.GET("/", service.ListAllRepos)

		repo_member_router := repo_router.Use(middleware.MemberCheck())
		{
			// repo_member_router.GET("/:repo_name/", service.GetRepo)
			repo_member_router.POST(" /", service.CreateRepo)
		}
	}

	kanban_router := project_router.Group("/:project_name/kanban")
	{
		kanban_router.GET("/", service.ListAllKanbans)

		kanban_member_router := kanban_router.Use(middleware.MemberCheck())
		{
			// kanban_member_router.GET("/", service.GetKanban)
			// kanban_member_router.POST("/", service.CreateKanban)
			kanban_member_router.PUT("/", service.AddWorkFlow)
			kanban_member_router.POST("/", service.AddTask)
			kanban_member_router.DELETE("/:workflow_name/:task_id/", service.DeleteTask)
		}
	}
}
