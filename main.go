package main

import (
	"github.com/Ivanhahanov/GoLibrary/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/books/:id", routes.HandleGetBook)
	r.GET("/books/", routes.HandleGetBooks)
	r.PUT("/books/", routes.HandleCreateBook)
	r.POST("/books/", routes.HandleUpdateBook)
	r.Run()
}
