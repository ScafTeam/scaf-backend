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

	// "github.com/ScafTeam/firebase-go-client/auth"
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"google.golang.org/api/iterator"
)

func ListAllProjects(c *gin.Context) {
	claims, err := middleware.AuthMiddleware.GetClaimsFromJWT(c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	userEmail := claims[middleware.IdentityKey].(string)
	log.Println("get projects")
	iter := database.Client.Collection("projects").
		Where("Members", "array-contains", userEmail).
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
	claims, err := middleware.AuthMiddleware.GetClaimsFromJWT(c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	userEmail := claims[middleware.IdentityKey].(string)
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
		"Members":  []string{userEmail},
		"Repos":    []model.Repo{},
		"DevTools": req["DevTools"],
		"DevMode":  req["DevMode"],
	}

	_, err = database.Client.
		Doc("projects/"+project_uuid).
		Set(context.Background(), project)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "project created",
	})
}

func DeleteProject(c *gin.Context) {

	project_id := c.Param("project_id")
	log.Println("delete project")
	_, err := database.Client.
		Collection("projects").
		Doc(project_id).
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

func AddMember(c *gin.Context) {
	project_id := c.Param("project_id")
	claims, err := middleware.AuthMiddleware.GetClaimsFromJWT(c)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	userEmail := claims[middleware.IdentityKey].(string)

	_, err = database.Client.
		Collection("projects").
		Doc(project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "Members",
				Value: firestore.ArrayUnion(userEmail),
			},
		})
	if err != nil {
		log.Printf("An error has occurred: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "member added",
	})
}
