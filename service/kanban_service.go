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

func createKanban(c *gin.Context) model.Kanban {
	project_id := c.Param("project_id")

	var workflows map[string]model.Workflow
	workflows = make(map[string]model.Workflow)

	todo_id := shortuuid.New()
	workflows[todo_id] = model.Workflow{
		Id:    todo_id,
		Name:  "Todo",
		Tasks: []model.Task{},
	}

	inProgress_id := shortuuid.New()
	workflows[inProgress_id] = model.Workflow{
		Id:    inProgress_id,
		Name:  "In Progress",
		Tasks: []model.Task{},
	}

	done_id := shortuuid.New()
	workflows[done_id] = model.Workflow{
		Id:    done_id,
		Name:  "Done",
		Tasks: []model.Task{},
	}

	kanban := model.Kanban{
		ProjectId: project_id,
		Workflows: workflows,
	}

	return kanban
}

func ListAllKanbans(c *gin.Context) {
	res := getProjectDetail(c)

	project_id := res.Ref.ID

	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Get(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	var kanban model.Kanban
	dsnap.DataTo(&kanban)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "List Kanban",
		"kanban":  kanban,
	})
}

func AddWorkFlow(c *gin.Context) {
	project_id := getProjectId(c)

	var req model.AddWorkFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	new_workflow_id := shortuuid.New()
	workflow := model.Workflow{
		Id:    new_workflow_id,
		Name:  req.WorkflowName,
		Tasks: []model.Task{},
	}

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Set(context.Background(), map[string]interface{}{
			"workflows." + new_workflow_id: workflow,
		}, firestore.MergeAll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Add workflow success",
	})
}

func AddTask(c *gin.Context) {
	project_id := getProjectId(c)
	task_id := shortuuid.New()

	var req model.AddTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check workflow_Id in firebase kanbans
	workflow := findWorkFlowById(c, project_id, req.WorkflowId)

	task := model.Task{
		Id:          task_id,
		Name:        req.Name,
		Description: req.Description,
	}

	workflow.Tasks = append(workflow.Tasks, task)

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Set(context.Background(), map[string]interface{}{
			"workflows." + req.WorkflowId: workflow,
		}, firestore.MergeAll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Add Todo Task Success",
	})
}

// TODO
func DeleteTask(c *gin.Context) {
	project_id := getProjectId(c)
	task_id := c.Param("task_id")

	req := model.DeleteTaskRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check workflow_name in firebase kanbans
	workflow := findWorkFlowById(c, project_id, req.WorkflowId)

	var new_tasks []model.Task
	for _, task := range workflow.Tasks {
		if task.Id != task_id {
			new_tasks = append(new_tasks, task)
		}
	}

	workflow.Tasks = new_tasks

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Set(context.Background(), map[string]interface{}{
			"workflows." + req.WorkflowId: workflow,
		}, firestore.MergeAll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Delete Task Success",
	})
}

func findWorkFlowById(c *gin.Context, project_id, workflow_Id string) model.Workflow {
	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Get(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		c.Abort()
		return model.Workflow{}
	}

	var kanban model.Kanban
	dsnap.DataTo(&kanban)

	workflow := kanban.Workflows[workflow_Id]

	if workflow.Id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Workflow not found",
		})
		c.Abort()
		return model.Workflow{}
	}
	return workflow
}
