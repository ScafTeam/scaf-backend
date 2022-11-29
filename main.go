package main

import (
	"backend/myauth"
	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	auth_router := server.Group("/auth")
	auth.Auth("AIzaSyAvQMZVhXbBZ61DdypPJG-zsg0NHnqKEBQ")
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	auth_router.POST("/login", myauth.UserLogin)
	auth_router.POST("/register", myauth.UserRegister)
	server.Run(":8000")
}
