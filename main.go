package main

import (
	"backend/database"
	"backend/router"
	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
	// "log"
)

func main() {
	server := gin.Default()
	auth.Auth("AIzaSyAvQMZVhXbBZ61DdypPJG-zsg0NHnqKEBQ")
	database.SetupFirebase()

	auth_router := server.Group("/auth")
	router.AddAuthRouter(auth_router)

	project_router := server.Group("/projects")
	router.AddProjectRouter(project_router)

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	// firebase.GetData()
	server.Run(":8000")
}
