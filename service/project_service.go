package service

import (
	"backend/database"
	"backend/middleware"
	"backend/model"
	"context"
	"net/http"

	"encoding/json"

	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"google.golang.org/api/iterator"
)

func ListAllProjects(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	user, _ := c.Get(middleware.IdentityKey)
	userEmail := user.(*model.ScafUser).Email
	log.Println("get all projects")
	iter := database.Client.Collection("all_project").
		Doc(userEmail).
		Collection("projects").
		Documents(context.Background())
	projects := make([]model.Project, 0)
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
		projects = append(projects, project)
	}
	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

func CreateProject(c *gin.Context) {
	log.Println("create project")
	user, _ := c.Get(middleware.IdentityKey)
	userEmail := user.(*model.ScafUser).Email
	req := make(map[string]interface{})
	c.BindJSON(&req)
	log.Println(req)

	project_uuid := shortuuid.New()
	// repo_uuid := uuid.New().String()
	project_create_on := time.Now().Format(time.RFC850)

	project := map[string]interface{}{
		"Id":       project_uuid,
		"Name":     req["Name"],
		"CreateOn": project_create_on,
		"Author":   userEmail,
		"Members":  []string{},
		"Repos":    []model.Repo{},
		"DevTools": req["DevTools"],
		"DevMode":  req["DevMode"],
	}

	_, err := database.Client.
		Doc("all_projects/"+userEmail+"/projects/"+project_uuid).
		Set(context.Background(), project)
	if err != nil {
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
	log.Println("add repo")
	user, _ := c.Get(middleware.IdentityKey)
	userEmail := user.(*model.ScafUser).Email
	req := make(map[string]interface{})
	c.BindJSON(&req)
	id := req["project_id"].(string)

	repo := model.Repo{
		Name: req["Name"].(string),
		Url:  req["Url"].(string),
	}

	dsnap, err := database.Client.
		Doc("all_projects/" + userEmail + "/projects/" + id).
		Get(context.Background())

	if err != nil {
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
		Doc("all_projects/" + userEmail + "/projects/" + id).
		Delete(context.Background())

	if err != nil {
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

	var mapData map[string]interface{}
	if err := json.Unmarshal(jsonStr, &mapData); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}

	_, err = database.Client.
		Doc("all_projects/"+userEmail+"/projects/"+id).
		Set(context.Background(), mapData)

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

func DeleteProject(c *gin.Context) {
	user, _ := c.Get(middleware.IdentityKey)
	userEmail := user.(*model.ScafUser).Email
	req := make(map[string]interface{})
	c.BindJSON(&req)

	log.Println("delete project")
	_, err := database.Client.
		Doc("all_projects/" + userEmail + "/projects/" + req["project_id"].(string)).
		Delete(context.Background())
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "project deleted",
	})
}

func ListAllRepos(c *gin.Context) {
	user, _ := c.Get(middleware.IdentityKey)
	userEmail := user.(*model.ScafUser).Email
	log.Println("list all repos")

	req := make(map[string]interface{})
	c.BindJSON(&req)

	dsnap, err := database.Client.
		Doc("all_projects/" + userEmail + "/projects/" + req["project_id"].(string)).
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
