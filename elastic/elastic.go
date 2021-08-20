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
	Highlight struct {
		AttachmentContent []string `json:"attachment.content"`
	} `json:"highlight"`
}

type Source struct {
	Path         string      `json:"path"`
	Year         string      `json:"year"`
	Author       string      `json:"author"`
	Description  string      `json:"description"`
	Publisher    string      `json:"publisher"`
	CreationDate time.Time   `json:"creation_date"`
	Title        string      `json:"title"`
	Slug         string      `json:"slug"`
	Tags         interface{} `json:"tags"`

	Attachment struct {
		ContentType   string `json:"content_type"`
		Language      string `json:"language"`
		ContentLength int    `json:"content_length"`
	} `json:"attachment"`
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
			"attachment": map[string]interface{}{
				"field":         "data",
				"indexed_chars": -1,
			},
		},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	req := esapi.IngestPutPipelineRequest{
		PipelineID: "attachment",
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
				"attachment.content": map[string]interface{}{
					"type":        "text",
					"analyzer":    analyzer,
					"term_vector": "with_positions_offsets",
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
