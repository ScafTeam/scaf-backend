package service

import (
	"backend/database"
	"backend/model"
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
)

func AddRepo(c *gin.Context) {
	var req model.AddRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": err.Error(),
		})
		return
	}

	project_id := getProjectId(c)

	repo := map[string]interface{}{
		"id":   shortuuid.New(),
		"name": req.Name,
		"url":  req.Url,
	}

	_, err := database.Client.
		Collection("projects").
		Doc(project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "repos",
				Value: firestore.ArrayUnion(repo),
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      repo["id"],
		"status":  "OK",
		"message": "Add Repo " + req.Name + " success",
	})
}

func UpdateRepo(c *gin.Context) {
	var req model.UpdateRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": err.Error(),
		})
		return
	}

	dsnap := getProjectDetail(c)
	var project model.Project
	dsnap.DataTo(&project)

	repos := project.Repos

	var hasRepo bool
	for i, repo := range repos {
		if repo.Id == req.Id {
			hasRepo = true
			if req.Name != "" {
				repos[i].Name = req.Name
			}
			if req.Url != "" {
				repos[i].Url = req.Url
			}
		}
	}

	if !hasRepo {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Repo not found",
		})
		return
	}

	_, err := database.Client.
		Doc("projects/"+dsnap.Ref.ID).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "repos",
				Value: repos,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Update repo " + req.Name + " success",
	})
}

func DeleteRepo(c *gin.Context) {
	var req model.DeleteRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "BadRequst",
			"message": err.Error(),
		})
		return
	}

	dsnap := getProjectDetail(c)
	var project model.Project
	dsnap.DataTo(&project)

	repos := project.Repos

	var hasRepo bool
	for i, repo := range repos {
		if repo.Id == req.Id {
			hasRepo = true
			repos = append(repos[:i], repos[i+1:]...)
			break
		}
	}

	if !hasRepo {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Repo not found",
		})
		return
	}

	_, err := database.Client.
		Doc("projects/"+dsnap.Ref.ID).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "repos",
				Value: processFirestoreData(repos),
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Delete repo success",
	})
}

func ListAllRepos(c *gin.Context) {
	dsnap := getProjectDetail(c)
	repos := dsnap.Data()["repos"].([]interface{})

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "List all repos",
		"repos":   repos,
	})
}

func processFirestoreData(data []model.Repo) []map[string]interface{} {
	res := []map[string]interface{}{}

	for _, v := range data {
		res = append(res, _processData(v))
	}

	return res
}
