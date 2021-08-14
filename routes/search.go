package routes

import (
	"fmt"
	"github.com/Ivanhahanov/GoLibrary/database"
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var sp SearchParams
	err := c.BindQuery(&sp)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"detail": "Invalid request params",
		})
	}
	search := database.MongoSearch{
		Title:       sp.Title,
		Author:      sp.Author,
		Description: sp.Description,
		Publisher:   sp.Publisher,
	}
	result := database.Search(&search)
	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}

type ContentSearch struct {
	Query             string `form:"q"`
	NumberOfFragments int    `form:"num_of_fragments"`
	Size              int    `form:"size"`
	Language          string `form:"lang"`
}

type Response struct {
	MongoId     string   `json:"mongo_id"`
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
	elasticResult := elastic.ContentSearch(index, cs.Query, cs.NumberOfFragments, cs.Size)
	var response []Response
	for _, result := range elasticResult {
		r := Response{
			MongoId: result.MongoID,
			Text:    result.Text,
		}
		objectId, err := primitive.ObjectIDFromHex(result.MongoID)
		if err != nil {
			log.Println(err)
		}
		book, err := database.GetBookByID(objectId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"detail": fmt.Sprintf("can't find book by id %s", result.MongoID),
			})
		}
		r.Title = book.Title
		r.Author = book.Author
		r.Publisher = book.Publisher
		r.Description = book.Description
		r.Tags = book.Tags
		r.Path = book.Path
		response = append(response, r)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": response,
	})
}
