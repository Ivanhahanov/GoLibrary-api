package main

import (
	"github.com/Ivanhahanov/GoLibrary/config"
	"github.com/Ivanhahanov/GoLibrary/database"
	"github.com/Ivanhahanov/GoLibrary/elastic"
	"github.com/Ivanhahanov/GoLibrary/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	var cfg config.Config
	cfg.LoadConfig("config.yml")
	elastic.InitConnection(&cfg)
	elastic.PipelineInit()
	elastic.IndexInit()
	database.InitConnection(&cfg)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.GET("/books/:id", routes.HandleGetBook)
	r.GET("/books/", routes.HandleGetBooks)
	r.PUT("/books/", routes.HandleUploadBook)
	r.POST("/books/", routes.HandleUpdateBook)
	r.DELETE("/books/:id", routes.HandleDeleteBook)

	r.GET("/search/", routes.HandleSearch)
	r.Run()
}
