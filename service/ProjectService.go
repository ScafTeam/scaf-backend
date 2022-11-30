package service

import (
	"backend/database"
	// "backend/model"
	"net/http"
	// "encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"time"
)

var projects []interface{}

func ListAllProjects(c *gin.Context) {
	log.Println("get all projects")
	dnaps, err := database.Client.Collection(User.Email).Doc("projects").Get(database.Ctx)
	if err != nil {
		log.Println(err)
		return
	}
	projects_json := dnaps.Data()["list"]
	projects = projects_json.([]interface{})
	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
	log.Println(projects)
}

func CreateProject(c *gin.Context) {
	log.Println("create project")
	project_uuid := uuid.New().String()
	repo_uuid := uuid.New().String()
	project_create_on := time.Now().Format(time.RFC850)
	_, err := database.Client.
		Doc("all_projects/"+User.Email+"/projects/"+project_uuid).
		Set(database.Ctx, map[string]interface{}{
			"Id":       project_uuid,
			"Name":     "test",
			"CreateOn": project_create_on,
			"Author":   User.Email,
			"Members":  []string{},
			"DevTools": []string{},
			"DevMode":  "waterfall",
		})
	_, err = database.Client.
		Doc("all_projects/"+User.Email+"/projects/"+project_uuid+"/Repos/"+repo_uuid).
		Set(database.Ctx, map[string]interface{}{
			"Id":   repo_uuid,
			"Name": "test",
			"Url":  "test",
		})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
}
