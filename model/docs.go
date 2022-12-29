package model

type Docs struct {
	ProjectId string         `json:"projectId"`
	Docs      map[string]Doc `json:"docs"`
}

type Doc struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AddDocRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateDocRequest struct {
	Id      string `json:"id" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type DeleteDocRequest struct {
	Id string `json:"id" binding:"required"`
}
