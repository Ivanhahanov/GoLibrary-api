package routes

import (
	"github.com/Ivanhahanov/GoLibrary/shortlinks"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Request struct {
	OriginalLink string `json:"original_link"`
}

func HandleShorter(c *gin.Context) {
	shortLink := c.Param("shortLink")
	log.Println(shortLink)
	originalLink, err := shortlinks.GetOriginalLink(shortLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	err = shortlinks.WriteVisit(shortLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.Redirect(http.StatusMovedPermanently, originalLink)
}

func HandleCreateShortLink(c *gin.Context) {
	var link Request
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(link)
	shortLink, err := shortlinks.CreateShortLink(link.OriginalLink)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Short link created!", "short_link": shortLink})
}

func HandleGetAllShortLinks(c *gin.Context) {
	list, err := shortlinks.GetAllDocuments()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": list})

}
