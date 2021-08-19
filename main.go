package main

import (
	"github.com/Ivanhahanov/GoLibrary/config"
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/Ivanhahanov/GoLibrary/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var cfg config.Config
	language := map[string]string {
		"books_ru": "russian",
		"books_en": "english",
	}
	cfg.LoadConfig("config.yml")
	elastic.InitConnection(&cfg)
	elastic.PipelineInit()
	for k, v := range language {
		elastic.IndexInit(k, v)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/books/:id", routes.HandleGetBook)
	r.GET("/books/", routes.HandleGetBooks)
	r.PUT("/books/", routes.HandleUploadBook)
	r.POST("/books/", routes.HandleUpdateBook)
	r.DELETE("/books/:id", routes.HandleDeleteBook)
	r.GET("/search/", routes.HandleSearch)
	r.GET("/content/search/", routes.HandleSearchContent)
	r.GET("/:shortLink", routes.HandleShorter)
	r.POST("/link/create/", routes.HandleCreateShortLink)
	r.GET("/link/", routes.HandleGetAllShortLinks)
	r.GET("/download/:id", routes.HandleDownload)
	r.Run(":80")
}
