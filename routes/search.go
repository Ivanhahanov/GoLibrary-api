package routes

import (
	"github.com/Ivanhahanov/GoLibrary/database"
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

func HandleSearchContent(c *gin.Context) {

}
