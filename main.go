package main

import (
	"backend/database"
	"backend/middleware"
	"backend/router"

	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	middleware.SetupAuthMiddleware(server)
	database.SetupFirebase()
	auth.Auth(database.Key)

	auth_router := server.Group("/")
	router.AddAuthRouter(auth_router)

	project_router := server.Group("user/:user_email/project")
	router.AddProjectRouter(project_router)

	repo_router := project_router.Group("/:project_name/repo")
	router.AddRepoRouter(repo_router)

	kanban_router := project_router.Group("/:project_name/kanban")
	router.AddKanbanRouter(kanban_router)

	server.Run(":8000")
}
