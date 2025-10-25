package elastic

import (
	"time"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/search-service/configs"
	"github.com/elastic/go-elasticsearch/v7"
)

var (
	ESClient *elasticsearch.Client
	Index    string
)

func Init(cfg *configs.Config) {
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.Elastic.URL},
	}

	var client *elasticsearch.Client
	var err error

	// Retry loop
	for i := 0; i < 10; i++ {
		client, err = elasticsearch.NewClient(esCfg)
		if err == nil {
			res, err := client.Info()
			if err == nil && res.StatusCode == 200 {
				res.Body.Close()
				break
			}
		}
		logger.Warnf("[elastic] waiting for Elasticsearch... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Errorf("[elastic] failed to connect to Elasticsearch: %v", err)
		return
	}

	ESClient = client
	Index = cfg.Elastic.Index
	logger.Infof("[elastic] connected to %s | index=%s", cfg.Elastic.URL, Index)
}
