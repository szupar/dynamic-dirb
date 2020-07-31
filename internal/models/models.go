package models

type Graph struct {
	Graph string `json:"graph"`
}

type Output struct {
	Name string   `json:"name"`
	Body []string `json:"output"`
}
