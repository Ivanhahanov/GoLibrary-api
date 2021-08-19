package routes

import (
	"github.com/gin-gonic/gin"
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
