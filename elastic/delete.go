package elastic

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
)

func Delete(mongoId string) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	req := esapi.DeleteRequest{Index: "books_ru", DocumentID: mongoId}
	if err != nil {
		log.Fatalf("Error deleting the document: %s", err)
	}
	res, err := req.Do(context.Background(), es)
	if err != nil{
		log.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}
