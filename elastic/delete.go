package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

func Delete(mongoId string) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"mongo_id": mongoId,
			},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := es.DeleteByQuery([]string{"books_ru"}, strings.NewReader(buf.String()))
	if err != nil {
		log.Fatalf("Error deleting the document: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}
