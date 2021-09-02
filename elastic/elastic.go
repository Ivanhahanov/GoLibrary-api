package elastic

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Ivanhahanov/GoLibrary/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
	"net/http"
	"time"
)

var cfg elasticsearch.Config

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
	Index     string  `json:"_index"`
	Type      string  `json:"_type"`
	ID        string  `json:"_id"`
	Score     float64 `json:"_score"`
	Source    Source  `json:"_source"`
	InnerHits struct {
		Attachments struct {
			Hits struct {
				Total struct {
					Value    int    `json:"value"`
					Relation string `json:"relation"`
				} `json:"total"`
				MaxScore float64 `json:"max_score"`
				Hits     []struct {
					Index  string `json:"_index"`
					Type   string `json:"_type"`
					ID     string `json:"_id"`
					Nested struct {
						Field  string `json:"field"`
						Offset int    `json:"offset"`
					} `json:"_nested"`
					Score     float64 `json:"_score"`
					Source    PageSource  `json:"_source"`
					Highlight struct {
						AttachmentsAttachmentContent []string `json:"attachments.attachment.content"`
					} `json:"highlight"`
				} `json:"hits"`
			} `json:"hits"`
		} `json:"attachments"`
	} `json:"inner_hits"`
}

type Source struct {
	Path         string        `json:"path"`
	Year         string        `json:"year"`
	Author       string        `json:"author"`
	Description  string        `json:"description"`
	Publisher    string        `json:"publisher"`
	CreationDate time.Time     `json:"creation_date"`
	Title        string        `json:"title"`
	Slug         string        `json:"slug"`
	Tags         interface{}   `json:"tags"`
	Text         []interface{} `json:"text"`
}
type PageSource struct {
	Data string `json:"data"`
	Page int `json:"page"`
}

func InitConnection(baseCfg *config.Config) {
	cfg = elasticsearch.Config{
		Addresses: []string{
			baseCfg.Elastic.Address,
		},
		Username: baseCfg.Elastic.Username,
		Password: baseCfg.Elastic.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func PipelineInit() {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"processors": []interface{}{map[string]interface{}{
			"foreach": map[string]interface{}{
				"field": "attachments",
				"processor": map[string]interface{}{
					"attachment": map[string]interface{}{
						"target_field":  "_ingest._value.attachment",
						"field":         "_ingest._value.data",
						"indexed_chars": -1,
					},
				},
			},
		},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	req := esapi.IngestPutPipelineRequest{
		PipelineID: "multiple_attachment",
		Body:       esutil.NewJSONReader(&query),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
func IndexInit(index_name string, analyzer string) {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"attachments": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"attachment": map[string]interface{}{
							"properties": map[string]interface{}{
								"content": map[string]interface{}{
									"type":        "text",
									"analyzer":    analyzer,
									"term_vector": "with_positions_offsets",
								},
							},
						},
					},
				},
			},
		},
		"settings": map[string]interface{}{
			"index": map[string]interface{}{
				"highlight": map[string]interface{}{
					"max_analyzed_offset": 1000000000,
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	req := esapi.IndicesCreateRequest{
		Index: index_name,
		Body:  esutil.NewJSONReader(&query),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}
