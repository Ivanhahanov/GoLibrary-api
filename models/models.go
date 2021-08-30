package models

type Book struct {
	Title        string        `json:"title"`
	Publisher    string        `json:"publisher"`
	Author       string        `json:"author"`
	Year         string        `json:"year"`
	Tags         []string      `json:"tags"`
	Description  string        `json:"description"`
	Slug         string        `json:"slug"`
	Path         string        `json:"path"`
	Data         []interface{} `json:"attachments"`
	CreationDate string        `json:"creation_date"`
}
