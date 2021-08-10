package elastic

import (
	"code.sajari.com/docconv"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
)

type Book struct {
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Tags   []string `json:"tags"`
	Data   string   `json:"data"`
}

func Put(filepath string) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	content, err := docconv.ConvertPath(filepath)
	if err != nil {
		log.Fatal(err)
	}

	str := base64.StdEncoding.EncodeToString([]byte(content.Body))

	book := Book{
		Title:  "test",
		Author: "test",
		Tags:   []string{"test", "AD"},
		Data:   str,
	}

	req := esapi.IndexRequest{
		Index:    "books_ru",
		Body:     esutil.NewJSONReader(&book),
		Pipeline: "attachment",
	}
	res, err := req.Do(context.Background(), es)
	defer res.Body.Close()
	fmt.Println(res)
}
