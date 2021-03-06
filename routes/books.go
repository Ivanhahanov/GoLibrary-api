package routes

import (
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/Ivanhahanov/GoLibrary/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

func HandleGetBooks(c *gin.Context) {
	var loadedBooks = elastic.GetAllInIndex("books_*")
	c.JSON(http.StatusOK, gin.H{"books": loadedBooks})
}

func HandleGetBook(c *gin.Context) {
	bookId := c.Param("id")
	var loadedBook = elastic.GetById(bookId)
	// TODO: return file by loadedBook.Path
	c.JSON(http.StatusOK, gin.H{"ID": loadedBook.Slug, "Path": loadedBook.Path})

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
	language := c.PostForm("language")
	index := "books_ru"
	switch language {
	case "en":
		index = "books_en"
	case "ru":
		index = "books_ru"
	}
	title := c.PostForm("title")
	author := c.PostForm("author")
	publisher := c.PostForm("publisher")
	description := c.PostForm("description")
	// year := c.PostForm("year")
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
	book.Publisher = publisher
	book.Description = description
	book.Year = "2021"
	book.Slug = slug.Make(book.Title) + "_" + book.Year
	book.CreationDate = time.Now().Format(time.RFC3339)

	elastic.Put(&book, index)
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
}

func HandleDeleteBook(c *gin.Context) {
	bookId := c.Param("id")
	elastic.Delete(bookId)
	c.JSON(http.StatusOK, gin.H{"book": bookId})
}

func HandleDownload(c *gin.Context) {
	bookId := c.Param("id")
	loadedBook := elastic.GetById(bookId)
	c.FileAttachment(loadedBook.Path, loadedBook.Slug+".pdf")
}
