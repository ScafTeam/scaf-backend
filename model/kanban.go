package model

type Kanban struct {
	ProjectId string     `json:"projectId"`
	Workflows []Workflow `json:"workflows"`
}

type Workflow struct {
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
