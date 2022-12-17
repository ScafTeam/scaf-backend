package model

type Kanban struct {
	ProjectId  string `json:"project_id"`
	Todo       []Task `json:"todo"`
	InProgress []Task `json:"in_progress"`
	Done       []Task `json:"done"`
}

type Task struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
