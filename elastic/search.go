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
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"nested": map[string]interface{}{
						"path": "attachments",
						"query": map[string]interface{}{
							"simple_query_string": map[string]interface{}{
								"query": searchString,
								"fields": []string{
									"attachments.attachment.content",
								},
								"default_operator": "and",
							},
						},
						"inner_hits": map[string]interface{}{
							"size":    4,
							"_source": false,
							"highlight": map[string]interface{}{
								"order":               "score",
								"pre_tags":            "<b>",
								"post_tags":           "</b>",
								"number_of_fragments": 1,
								"fragment_size":       25,
								"fields": map[string]interface{}{
									"attachments.attachment.content": map[string]interface{}{},
								},
							},
						},
					},
				},
				"should": map[string]interface{}{
					"simple_query_string": map[string]interface{}{
						"query": searchString,
						"fields": []string{
							"title", "description",
						},
						"default_operator": "and",
					},
				},
			},
		},
		"_source": map[string]interface{}{
			"excludes": []string{
				"attachments",
			},
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
