package main

import (
	"github.com/Ivanhahanov/GoLibrary/config"
	"github.com/Ivanhahanov/GoLibrary/database"
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/Ivanhahanov/GoLibrary/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	var cfg config.Config
	cfg.LoadConfig("config.yml")
	elastic.InitConnection(&cfg)
	database.InitConnection(&cfg)
	r := gin.Default()
	r.GET("/books/:id", routes.HandleGetBook)
	r.GET("/books/", routes.HandleGetBooks)
	r.PUT("/books/", routes.HandleUploadBook)
	r.POST("/books/", routes.HandleUpdateBook)
	r.Run()
}
