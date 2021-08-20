package elastic

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)


func ContentSearch(index string, searchString string, numberOfFragments int, fragmentSize int) (output []SearchItem) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"simple_query_string": map[string]interface{}{
				"query":            searchString,
				"fields":           []string{"attachment.content"},
				"default_operator": "and",
			},
		},
		"_source": map[string]interface{}{
			"excludes": []string{"data", "attachment.content"},
		},
		"highlight": map[string]interface{}{
			"order":               "score",
			"number_of_fragments": numberOfFragments,
			"fragment_size":       fragmentSize,
			"pre_tags":            "<b>",
			"post_tags":           "</b>",
			"fields": map[string]interface{}{
				"attachment.content": map[string]interface{}{},
			},
			"type": "fvh",
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.Search(es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty())
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var r SearchResult

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	if r.Hits.Total.Value > 0 {

		// Print the response status, number of results, and request duration.
		log.Printf(
			"[%s] %d hits; took: %dms",
			res.Status(),
			int(r.Hits.Total.Value),
			int(r.Took),
		)
		// Print the ID and document source for each hit.
		for _, hit := range r.Hits.Hits {
			output = append(output, hit)
		}
	}
	return
}
