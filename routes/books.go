package routes

import (
	"github.com/Ivanhahanov/GoLibrary/database"
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/Ivanhahanov/GoLibrary/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
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

type FileJson struct {
	Name     string                `form:"name" binding:"required"`
	Author   string                `form:"author"`
	fileData *multipart.FileHeader `form:"file"`
}

func HandleUploadBook(c *gin.Context) {
	file, err := c.FormFile("file")

	// The file cannot be received.
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}
	title := c.Request.PostFormValue("title")
	author := c.Request.PostFormValue("author")

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.New().String() + extension

	// The file is received, so let's save it
	newFilePath := "/etc/golibrary/books/" + newFileName
	if err := c.SaveUploadedFile(file, newFilePath); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}
	// TODO: verify fields
	var book models.Book
	book.Author = author
	book.Title = title
	book.Path = newFilePath
	_, err = database.Create(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": book.Tags})
		return
	}

	elastic.Put(newFilePath)

	c.JSON(http.StatusOK, gin.H{
		"message": "Your file has been successfully uploaded.",
	})
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
