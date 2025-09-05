package elastic

import (
	"log"

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

    client, err := elasticsearch.NewClient(esCfg)
    if err != nil {
        logger.Errorf("[elastic] failed to create client: %v", err)
        return
    }

    // Проверка соединения
    res, err := client.Info()
    if err != nil {
        logger.Errorf("[elastic] failed to connect to Elasticsearch: %v", err)
        return
    }
    defer res.Body.Close()

    ESClient = client
    Index = cfg.Elastic.Index
    logger.Infof("[elastic] connected to %s | index=%s", cfg.Elastic.URL, Index)
}
