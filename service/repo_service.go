package service

import (
	"backend/database"
	"backend/model"
	"context"
	"net/http"

	"encoding/json"

	"log"

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

	dsnap, err := database.Client.
		Doc("projects/" + project_id).
		Get(context.Background())

	if err != nil {
		log.Printf("An get project error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	jsonStr, err := json.Marshal(dsnap.Data())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	var project model.Project
	if err := json.Unmarshal(jsonStr, &project); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	project.Repos = append(project.Repos, repo)

	jsonStr, err = json.Marshal(project)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	var mapData map[string]interface{}
	if err := json.Unmarshal(jsonStr, &mapData); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	_, err = database.Client.
		Doc("projects/"+project_id).
		Set(context.Background(), mapData)

	if err != nil {
		log.Printf("An add project error has occurred: %s", err)
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
