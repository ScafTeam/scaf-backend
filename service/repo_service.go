package service

import (
	"backend/database"
	"backend/model"
	"context"
	"net/http"

	"log"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func AddRepo(c *gin.Context) {
	log.Println("add repo")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	project_id := c.Param("project_id")

	repo := model.Repo{
		Name: req["Name"].(string),
		Url:  req["Url"].(string),
	}

	_, err := database.Client.
		Collection("projects").
		Doc(project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "Repos",
				Value: firestore.ArrayUnion(repo),
			},
		})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "add repo",
	})
}

func ListAllRepos(c *gin.Context) {
	project_id := c.Param("project_id")
	log.Println("list all repos")

	req := make(map[string]interface{})
	c.BindJSON(&req)

	dsnap, err := database.Client.
		Doc("projects/" + project_id).
		Get(context.Background())

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	repos := dsnap.Data()["Repos"].([]interface{})

	c.JSON(http.StatusOK, gin.H{
		"repos": repos,
	})
}
