package main

import (
	"backend/database"
	"backend/router"
	"backend/service"

	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
	// "log"
)

func main() {
	server := gin.Default()
	auth.Auth("AIzaSyAvQMZVhXbBZ61DdypPJG-zsg0NHnqKEBQ")
	database.SetupFirebase()

	auth_router := server.Group("/")
	router.AddAuthRouter(auth_router)

	project_router := server.Group("/project")
	router.AddProjectRouter(project_router)

	server.GET("/projects", service.ListAllProjects)
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	// firebase.GetData()
	server.Run(":8000")
}
