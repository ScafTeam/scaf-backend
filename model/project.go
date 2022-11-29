package model

type Project struct {
	Name     string   `json:"Name"`
	Author   string   `json:"Author"`
	CreateOn string   `json:"CreateOn"`
	Repos    []string `json:"Repos"`
	Members  []string `json:"Members"`
	DevTools []string `json:"DevTools"`
	DevMode  string   `json:"DevMode"`
}
