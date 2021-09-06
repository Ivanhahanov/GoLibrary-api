package routes

import (
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type SearchParams struct {
	Language    string `form:"lang"`
	Type        string `form:"type"`
	Title       string `form:"title"`
	Author      string `form:"author"`
	Description string `form:"description"`
	Publisher   string `form:"publisher"`
}

type Response struct {
	Year         string      `json:"year"`
	Author       string      `json:"author"`
	Description  string      `json:"description"`
	Publisher    string      `json:"publisher"`
	CreationDate time.Time   `json:"creation_date"`
	Title        string      `json:"title"`
	Slug         string      `json:"slug"`
	Tags         interface{} `json:"tags"`
	Text         interface{} `json:"text"`
}

type AutocompleteRequest struct {
	Query			string `form:"q"`
	Language		string `form:"lang"`
}

func HandleAutocompleteTitle(c *gin.Context) {
	var at AutocompleteRequest
	err := c.BindQuery(&at)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Invalid request params",
		})
	}
	index := "books_ru"
	switch at.Language {
	case "en":
		index = "books_en"
	case "ru":
		index = "books_ru"
	}
	response := elastic.AutocompleteTitle(index, at.Query)
	c.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}

type ContentSearch struct {
	Query             string `form:"q"`
	NumberOfFragments int    `form:"num_of_fragments"`
	Size              int    `form:"size"`
	Language          string `form:"lang"`
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
	results := elastic.ContentSearch(index, cs.Query, cs.NumberOfFragments, cs.Size)
	var response []*Response
	for _, result := range results {
		response = append(response, &Response{
			Year:         result.Source.Year,
			Author:       result.Source.Author,
			Description:  result.Source.Description,
			Publisher:    result.Source.Publisher,
			CreationDate: result.Source.CreationDate,
			Title:        result.Source.Title,
			Slug:         result.Source.Slug,
			Tags:         result.Source.Tags,
			Text:         result.Source.Text,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"result": response,
	})

}

type GetPagesRequest struct {
	Language string `form:"lang"`
	Slug     string `form:"slug"`
	Page     int    `form:"page"`
}

func HandleGetPages(c *gin.Context) {
	var pr GetPagesRequest
	err := c.BindQuery(&pr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Invalid request params",
		})
	}
	index := "books_ru"
	switch pr.Language {
	case "en":
		index = "books_en"
	case "ru":
		index = "books_ru"
	}
	pages := elastic.GetPages(
		index,
		pr.Slug,
		pr.Page,
	)
	c.FileAttachment(pages,
		strings.Split(pages, "/")[1])
	os.Remove(pages)
}
