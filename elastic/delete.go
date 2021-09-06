package elastic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

func Delete(id string) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"slug": id,
			},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := es.DeleteByQuery([]string{"books_*"}, strings.NewReader(buf.String()))
	if err != nil {
		log.Fatalf("Error deleting the document: %s", err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}
