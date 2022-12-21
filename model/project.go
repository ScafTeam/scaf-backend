package model

import "strings"

type Project struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Author   string   `json:"author"`
	CreateOn string   `json:"createOn"`
	Repos    []Repo   `json:"repos"`
	Members  []string `json:"members"`
	DevTools []string `json:"devTools"`
	DevMode  string   `json:"devMode"`
}

type Repo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type CreateProjectRequest struct {
	Name     string   `json:"name" binding:"required"`
	DevTools []string `json:"devTools"`
	DevMode  string   `json:"devMode" binding:"required"`
}

type AddMemberRequest struct {
	Email string `json:"email" binding:"email"`
}

type DeleteMemberRequest struct {
	Email string `json:"email" binding:"email"`
}

type AddRepoRequest struct {
	Name string `json:"name" binding:"required"`
	Url  string `json:"url" binding:"required"`
}

type UpdateRepoRequest struct {
	Id   string `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Url  string `json:"url" binding:"required"`
}

type DeleteRepoRequest struct {
	Id string `json:"id" binding:"required"`
}

func CheckProjectName(name string) bool {
	// check project name without invalid characters
	return !strings.ContainsAny(name, "/\\?%*:|\"<>")
}
