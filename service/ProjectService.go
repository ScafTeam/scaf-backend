package service

import (
	"backend/database"
	"backend/model"
	"net/http"

	"encoding/json"

	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

var Projects []interface{}

func ListAllProjects(c *gin.Context) {
	if User.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}
	log.Println("get all projects")
	iter := database.Client.Collection("all_projects").Doc(User.Email).Collection("projects").Documents(database.Ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
			})
			return
		}
		log.Println(doc.Data())
		jsonStr, err := json.Marshal(doc.Data())
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
		Projects = append(Projects, project)
	}
	c.JSON(http.StatusOK, gin.H{
		"projects": Projects,
	})
}

func CreateProject(c *gin.Context) {
	if User.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}
	log.Println("create project")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	log.Println(req)

	project_uuid := uuid.New().String()
	// repo_uuid := uuid.New().String()
	project_create_on := time.Now().Format(time.RFC850)

	project := map[string]interface{}{
		"Id":       project_uuid,
		"Name":     req["Name"],
		"CreateOn": project_create_on,
		"Author":   User.Email,
		"Members":  []string{},
		"Repos":    []model.Repo{},
		"DevTools": req["DevTools"],
		"DevMode":  req["DevMode"],
	}

	_, err := database.Client.
		Doc("all_projects/"+User.Email+"/projects/"+project_uuid).
		Set(database.Ctx, project)
	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "project created",
	})
}

func AddRepo(c *gin.Context) {
	if User.Email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}
	id := c.Param("id")
	log.Println("add repo")
	req := make(map[string]interface{})
	c.BindJSON(&req)
	log.Println(req)

	repo := model.Repo{
		Name: req["Name"].(string),
		Url:  req["Url"].(string),
	}

	dsnap, err := database.Client.
		Doc("all_projects/" + User.Email + "/projects/" + id).
		Get(database.Ctx)

	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
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

	_, err = database.Client.
		Doc("all_projects/" + User.Email + "/projects/" + id).
		Delete(database.Ctx)

	if err != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	project.Repos = append(project.Repos, repo)

	jsonStr, err = json.Marshal(project)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	// Convert struct
	var mapData map[string]interface{}
	if err := json.Unmarshal(jsonStr, &mapData); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	_, err = database.Client.
		Doc("all_projects/"+User.Email+"/projects/"+id).
		Set(database.Ctx, mapData)

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "add repo",
	})
	// _, err := database.Client.
	// 	Doc("all_projects/" + User.Email + "/projects/" + id).
	// 	Set(database.Ctx)
}
