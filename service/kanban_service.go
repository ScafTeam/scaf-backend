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

func createKanban(c *gin.Context, project_id string) map[string]interface{} {
	workflows := make(map[string]model.Workflow)

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

	kanban := map[string]interface{}{
		"projectId": project_id,
		"workflows": workflows,
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
		Update(context.Background(), []firestore.Update{
			{
				Path:  "workflows." + new_workflow_id,
				Value: workflow,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      new_workflow_id,
		"status":  "Created",
		"message": "Add workflow " + req.WorkflowName + " success",
	})
}

func UpdateWorkFlow(c *gin.Context) {
	var req model.UpdateWorkFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	project_id := getProjectId(c)
	workflow, ok := findWorkFlowById(c, project_id, req.Id)

	if !ok {
		return
	}

	workflow.Name = req.Name

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "workflows." + req.Id,
				Value: workflow,
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
		"message": "Update " + workflow.Name + " Success",
	})
}

func DeleteWorkFlow(c *gin.Context) {
	var req model.DeleteWorkFlowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	project_id := getProjectId(c)
	workflow, ok := findWorkFlowById(c, project_id, req.Id)

	if !ok {
		return
	}

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "workflows." + req.Id,
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
		"message": "Delete " + workflow.Name + " Success",
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
	workflow, ok := findWorkFlowById(c, project_id, req.WorkflowId)

	if !ok {
		return
	}

	task := model.Task{
		Id:          task_id,
		Name:        req.Name,
		Description: req.Description,
	}

	workflow.Tasks = append(workflow.Tasks, task)

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "workflows." + req.WorkflowId,
				Value: workflow,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Add " + workflow.Name + " Task Success",
	})
}

func UpdateTask(c *gin.Context) {
	project_id := getProjectId(c)

	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check workflow_Id in firebase kanbans
	workflow, ok := findWorkFlowById(c, project_id, req.WorkflowId)

	if !ok {
		return
	}

	var hasTask bool
	var new_tasks []model.Task
	for _, t := range workflow.Tasks {
		if t.Id == req.TaskId {
			hasTask = true
			t.Name = req.Name
			t.Description = req.Description
		}
		new_tasks = append(new_tasks, t)
	}

	if !hasTask {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Task not found",
		})
		return
	}

	workflow.Tasks = new_tasks

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "workflows." + req.WorkflowId,
				Value: workflow,
			},
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Update " + workflow.Name + " Task Success",
	})
}

func MoveTask(c *gin.Context) {
	project_id := getProjectId(c)

	var req model.MoveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check workflow_Id in firebase kanbans
	workflow, ok := findWorkFlowById(c, project_id, req.WorkflowId)

	if !ok {
		return
	}

	var task *model.Task
	var new_tasks []model.Task
	for _, t := range workflow.Tasks {
		if t.Id == req.TaskId {
			task = &t
		} else {
			new_tasks = append(new_tasks, t)
		}
	}

	if task == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Task not found",
		})
		return
	}

	// check new_workflow_Id in firebase kanbans
	new_workflow, ok := findWorkFlowById(c, project_id, req.WorkflowId)

	if !ok {
		return
	}

	new_workflow.Tasks = append(new_workflow.Tasks, *task)

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path: "workflows." + req.WorkflowId,
				Value: model.Workflow{
					Id:    workflow.Id,
					Name:  workflow.Name,
					Tasks: new_tasks,
				},
			},
			{
				Path:  "workflows." + req.NewWorkflowId,
				Value: new_workflow,
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
		"message": "Move Task Success",
	})
}

func DeleteTask(c *gin.Context) {
	project_id := getProjectId(c)

	req := model.DeleteTaskRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": err.Error(),
		})
		return
	}

	// check workflow_name in firebase kanbans
	workflow, ok := findWorkFlowById(c, project_id, req.WorkflowId)

	if !ok {
		return
	}

	var delete_task *model.Task
	var new_tasks []model.Task
	for _, task := range workflow.Tasks {
		if task.Id != req.TaskId {
			new_tasks = append(new_tasks, task)
		} else {
			delete_task = &task
		}
	}

	if delete_task == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Task not found",
		})
		return
	}

	workflow.Tasks = new_tasks

	_, err := database.Client.
		Doc("kanbans/"+project_id).
		Update(context.Background(), []firestore.Update{
			{
				Path:  "workflows." + req.WorkflowId,
				Value: workflow,
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
		"message": "Delete Task Success",
	})
}

func findWorkFlowById(c *gin.Context, project_id, workflow_Id string) (model.Workflow, bool) {
	dsnap, err := database.Client.
		Doc("kanbans/" + project_id).
		Get(context.Background())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": err.Error(),
		})
		c.Abort()
		return model.Workflow{}, false
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
		return model.Workflow{}, false
	}
	return workflow, true
}
