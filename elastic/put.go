package elastic

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Ivanhahanov/GoLibrary/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/gosimple/slug"
	"io/ioutil"
	"log"
	"os"
)


func Put(book *models.Book) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	f, _ := os.Open(book.Path)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	//if err != nil {
	//	log.Fatal(err)
	//}

	str := base64.StdEncoding.EncodeToString(content)

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

