package service

import (
	"backend/database"
	"backend/model"
	// "net/http"

	// "encoding/json"

	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

var projects []interface{}

func ListAllProjects(c *gin.Context) {
	log.Println("get all projects")
	iter := database.Client.Collection("all_projects").Doc(User.Email).Collection("/projects").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(doc.Data())
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"projects": projects,
	// })
}

func CreateProject(c *gin.Context) {
	log.Println("create project")
	json := make(map[string]interface{})
	c.BindJSON(&json)
	project_uuid := uuid.New().String()
	// repo_uuid := uuid.New().String()
	project_create_on := time.Now().Format(time.RFC850)
	_, err := database.Client.
		Doc("all_projects/"+User.Email+"/projects/"+project_uuid).
		Set(database.Ctx, map[string]interface{}{
			"Id":       project_uuid,
			"Name":     "test",
			"CreateOn": project_create_on,
			"Author":   User.Email,
			"Members":  []string{},
			"Repos": []model.Repo{
				{
					Name: "test",
					Url:  "asdfasdf",
				},
			},
			"DevTools": []string{},
			"DevMode":  "waterfall",
		})
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
	}
}
