package models

type Book struct {
	Title       string             `json:"title"`
	Publisher   string             `json:"publisher"`
	Author      string             `json:"author"`
	Tags        []string           `json:"tags"`
	Description string             `json:"description"`
	Path        string             `json:"path"`
	Data    string `json:"data"`
}
