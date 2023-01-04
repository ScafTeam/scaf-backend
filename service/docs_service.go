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

func createDocs(c *gin.Context, project_id string) map[string]interface{} {
	docs := make(map[string]model.Doc)

	docsMap := map[string]interface{}{
		"projectId": project_id,
		"docs":      docs,
	}

	return docsMap
}

func ListAllDocs(c *gin.Context) {
	projectID := getProjectId(c)

	dsnap, err := database.Client.
		Doc("docs/" + projectID).
		Get(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	var docs model.Docs
	dsnap.DataTo(&docs)

	var docsList []model.Doc
	for _, doc := range docs.Docs {
		docsList = append(docsList, doc)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Get all docs successfully",
		"docs": map[string]interface{}{
			"ProjecId": docs.ProjectId,
			"docs":     docsList,
		},
	})
}

func AddDoc(c *gin.Context) {
	var req model.AddDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	projectID := getProjectId(c)
	doc := model.Doc{
		Id:      shortuuid.New(),
		Title:   req.Title,
		Content: req.Content,
	}

	_, err := database.Client.
		Doc("docs/"+projectID).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "docs." + doc.Id,
				Value: doc,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Add doc " + req.Title + " successfully",
		"id":      doc.Id,
	})
}

func UpdateDoc(c *gin.Context) {
	var req model.UpdateDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	if !CheckDocExistByDocId(c, req.Id) {
		return
	}

	projectID := getProjectId(c)
	doc := model.Doc(req)

	_, err := database.Client.
		Doc("docs/"+projectID).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "docs." + doc.Id,
				Value: doc,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Update doc " + req.Title + " successfully",
	})
}

func DeleteDoc(c *gin.Context) {
	var req model.DeleteDocRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	if !CheckDocExistByDocId(c, req.Id) {
		return
	}

	projectID := getProjectId(c)

	_, err := database.Client.
		Doc("docs/"+projectID).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "docs." + req.Id,
				Value: firestore.Delete,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Delete doc " + req.Id + " successfully",
	})
}

func CheckDocExistByDocId(c *gin.Context, docId string) bool {
	projectID := getProjectId(c)

	dsnap, err := database.Client.
		Doc("docs/" + projectID).
		Get(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return false
	}

	var docs model.Docs
	dsnap.DataTo(&docs)

	if _, ok := docs.Docs[docId]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Doc not found",
		})
		return false
	}

	return true
}
