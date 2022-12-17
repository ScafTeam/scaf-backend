package main

import (
	"backend/database"
	"backend/middleware"
	"backend/router"
	"backend/service"

	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
	// "log"
)

func main() {
	server := gin.Default()
	auth.Auth("AIzaSyAvQMZVhXbBZ61DdypPJG-zsg0NHnqKEBQ")

	middleware.SetupAuthMiddleware(server)
	database.SetupFirebase()

	auth_router := server.Group("/")
	router.AddAuthRouter(auth_router)

	project_router := server.Group("/:user_email/project")
	router.AddProjectRouter(project_router)

	repo_router := project_router.Group("/:project_id/repo")
	router.AddRepoRouter(repo_router)

	server.GET("/projects", service.ListAllProjects)
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	// firebase.GetData()
	// p := middleware.AuthMiddleware.Authenticator
	server.Run(":8000")
}
