package elastic

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func GetById(documentId string) Source {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": documentId,
			},
		},
		"_source": map[string]interface{}{
			"excludes": []string{"data", "attachment"},
		},
	}
	searchResult := Search(query, "books_*")
	if len(searchResult) == 1{
		return searchResult[0]
	}
	return Source{}
}

func GetAllInIndex(index string) []Source {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"_source": map[string]interface{}{
			"excludes": []string{"data", "attachment"},
		},
	}
	return Search(query, index)
}

func Search(query map[string]interface{}, index string) (output []Source) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := es.Search(es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty())
	if err != nil {
		log.Fatalf("Error deleting the document: %s", err)
	}
	defer res.Body.Close()
	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	if r.Hits.Total.Value > 0 {
		log.Printf(
			"[%s] %d hits; took: %dms",
			res.Status(),
			r.Hits.Total.Value,
			r.Took,
		)
		for _, hit := range r.Hits.Hits {
			output = append(output, hit.Source)
		}
	}
	return
}
