package elastic

import (
	"code.sajari.com/docconv"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Ivanhahanov/GoLibrary/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/gosimple/slug"
	"log"
)


func Put(book *models.Book) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	content, err := docconv.ConvertPath(book.Path)
	if err != nil {
		log.Fatal(err)
	}

	str := base64.StdEncoding.EncodeToString([]byte(content.Body))

	book.Data = str
	req := esapi.IndexRequest{
		Index:    "books_ru",
		DocumentID: slug.Make(book.Title) + "_" + book.Year,
		Body:     esutil.NewJSONReader(&book),
		Pipeline: "attachment",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil{
		log.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}

