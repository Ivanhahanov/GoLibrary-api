package routes

import (
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type SearchParams struct {
	Language    string `form:"lang"`
	Type        string `form:"type"`
	Title       string `form:"title"`
	Author      string `form:"author"`
	Description string `form:"description"`
	Publisher   string `form:"publisher"`
}

func HandleSearch(c *gin.Context) {
}

type ContentSearch struct {
	Query             string `form:"q"`
	NumberOfFragments int    `form:"num_of_fragments"`
	Size              int    `form:"size"`
	Language          string `form:"lang"`
}

type Response struct {
	Title       string   `json:"title"`
	Publisher   string   `json:"publisher"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	Text        []string `json:"text"`
}

func HandleSearchContent(c *gin.Context) {
	var cs ContentSearch
	err := c.BindQuery(&cs)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Invalid request params",
		})
	}
	index := "books_ru"
	switch cs.Language {
	case "en":
		index = "books_en"
	case "ru":
		index = "books_ru"
	}
	log.Println(cs.Size)
	if cs.Size < 18 {
		cs.Size = 25
	}
	if cs.NumberOfFragments < 1 {
		cs.NumberOfFragments = 3
	}
	log.Println(cs.Query)
	_ = elastic.ContentSearch(index, cs.Query, cs.NumberOfFragments, cs.Size)
	var response []Response

	c.JSON(http.StatusOK, gin.H{
		"result": response,
	})

}
