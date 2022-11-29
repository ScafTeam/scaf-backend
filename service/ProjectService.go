package service

import (
	"backend/database"
	"backend/model"
	"fmt"
	// "github.com/ScafTeam/firebase-go-client/auth"
	"github.com/gin-gonic/gin"
)

var projects = []model.Project{}

func ListAllProjects(c *gin.Context) {
	projects, err := database.GetData(User.Email, "projects")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(projects)
}
