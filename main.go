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
	auth.Auth(database.Key)

	middleware.SetupAuthMiddleware(server)
	database.SetupFirebase()

	auth_router := server.Group("/")
	router.AddAuthRouter(auth_router)

	project_router := server.Group("/:user_email/project")
	router.AddProjectRouter(project_router)

	repo_router := project_router.Group("/:project_id/repo")
	router.AddRepoRouter(repo_router)

	server.Run(":8000")
}
