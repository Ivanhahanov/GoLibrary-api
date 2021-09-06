package elastic

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	pdf "github.com/pdfcpu/pdfcpu/pkg/api"
	"io/ioutil"
	"log"
	"os"
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
	if len(searchResult) == 1 {
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

func GetPages(index string, slug string, centralPage int) string {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": map[string]interface{}{
					"nested": map[string]interface{}{
						"path": "attachments",
						"query": map[string]interface{}{
							"terms": map[string]interface{}{
								"attachments.page": []int{
									centralPage - 1,
									centralPage,
									centralPage + 1,
								},
							},
						},
						"inner_hits": map[string]interface{}{
							"_source": map[string]interface{}{
								"includes": []string{"attachments.data"},
							},
						},
					},
				},
				"should": map[string]interface{}{
					"match": map[string]interface{}{
						"slug": slug,
					},
				},
			},
		},
		"_source": false,
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
		var files []string
		for _, hit := range r.Hits.Hits {
			for _, inner := range hit.InnerHits.Attachments.Hits.Hits {
				file, _ := base64.StdEncoding.DecodeString(inner.Source.Data)
				tmpFile, err := ioutil.TempFile("", "*.pdf")
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(len(file))
				if _, err := tmpFile.Write(file); err != nil {
					fmt.Println(err)
				}

				files = append(files, tmpFile.Name())
			}
		}
		fmt.Println(files)
		_ = pdf.MergeAppendFile(files, "/tmp/part.pdf", nil)
		for _, path := range files{
			os.Remove(path)
		}
		return "/tmp/part.pdf"
	}
	return ""
}

func AutocompleteTitle(index string, text string) []string {
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"autocomplete": map[string]interface{}{
				"prefix": text,
				"completion": map[string]interface{}{
					"field": "title",
				},
			},
		},
		"_source": false,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := es.Search(es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithPretty())
	if err != nil {
		log.Fatalf("Error deleting the document: %s", err)
	}
	defer res.Body.Close()

	var r Autocomplete
	var titles []string
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	for _, opt := range r.Suggest.Autocomplete {
		for _, text := range opt.Options {
			titles = append(titles, text.Text)
		}
	}
	return titles
}