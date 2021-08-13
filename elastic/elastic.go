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
)

var cfg elasticsearch.Config

func InitConnection(baseCfg *config.Config) {
	cfg = elasticsearch.Config{
		Addresses: []string{
			baseCfg.Elastic.Address,
		},
		// TODO: remove plain creds
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
func IndexInit() {
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
					"analyzer":    "russian",
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
	req:= esapi.IndicesCreateRequest{
		Index: "books",
		Body: esutil.NewJSONReader(&query),
	}
	res, err := req.Do(context.Background(), es)
	if err != nil{
		log.Println(err)
	}
	defer res.Body.Close()
	fmt.Println(res)
}
