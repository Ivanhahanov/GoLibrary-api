package routes

import (
	"github.com/Ivanhahanov/GoLibrary/database"
	"github.com/Ivanhahanov/GoLibrary/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func HandleGetBooks(c *gin.Context) {
	var loadedBooks, err = database.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": loadedBooks})
}

func HandleGetBook(c *gin.Context) {
	var book models.Book
	if err := c.BindUri(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	var loadedBook, err = database.GetBookByID(book.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	// TODO: return file by loadedBook.Path
	c.JSON(http.StatusOK, gin.H{"ID": loadedBook.ID, "Path": loadedBook.Path})
}

func HandleCreateBook(c *gin.Context) {
	// TODO: add file upload
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	// TODO: verify fields
	id, err := database.Create(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func HandleUpdateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	// TODO: check if book exists
	savedBook, err := database.Update(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": savedBook})
}
