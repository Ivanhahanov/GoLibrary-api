package elastic

import (
	"crypto/tls"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
)

var cfg = elasticsearch.Config{
	Addresses: []string{
		"https://localhost:9200",
	},
	// TODO: remove plain creds
	Username: "admin",
	Password: "admin",
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}
