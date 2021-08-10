package elastic

import (
	"crypto/tls"
	"github.com/Ivanhahanov/GoLibrary/config"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"

)


var cfg elasticsearch.Config

func InitConnection(baseCfg *config.Config){
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
