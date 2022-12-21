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

	router.Init(server)

	server.Run(":8000")
}
