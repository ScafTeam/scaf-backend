package main

import (
	"backend/database"
	"backend/middleware"
	"backend/router"

	"github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
    "bufio"
    "os"
)

func main() {
	server := gin.Default()
    key := readKey()
	auth.Auth(key)

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

func readKey() string {
  file, err := os.Open("database/key.txt")
  if err != nil {
    panic(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  scanner.Scan()
  return scanner.Text()
}
