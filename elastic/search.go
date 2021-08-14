package elastic

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type SearchResult struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64      `json:"max_score"`
		Hits     []SearchItem `json:"hits"`
	} `json:"hits"`
}

type SearchItem struct {
	Index  string  `json:"_index"`
	Type   string  `json:"_type"`
	ID     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source struct {
		MongoId string `json:"mongo_id"`
	} `json:"_source"`
	Highlight struct {
		AttachmentContent []string `json:"attachment.content"`
	} `json:"highlight"`
}

type OutputSearchResult struct {
	MongoID string   `json:"mongo_id"`
	Text    []string `json:"text"`
}

func ContentSearch(index string, searchString string, numberOfFragments int, fragmentSize int) (output []*OutputSearchResult) {
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
			"includes": "mongo_id",
		},
		"highlight": map[string]interface{}{
			"order":               "score",
			"number_of_fragments": numberOfFragments,
			"fragment_size":       fragmentSize,
			"pre_tags": "<b>",
			"post_tags": "</b>",
			"fields": map[string]interface{}{
				"attachment.content": map[string]interface{}{},
			},
			"type": "fvh",
		},
	}
	log.Println(query)
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
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r SearchResult

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Println(r)

	if r.Hits.Total.Value > 0 {

		for _, hit := range r.Hits.Hits {
			output = append(output, &OutputSearchResult{
				MongoID: hit.Source.MongoId,
				Text:    hit.Highlight.AttachmentContent,
			})
		}
	}
	return
}
