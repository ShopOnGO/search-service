package elastic

import (
    "log"

    "github.com/elastic/go-elasticsearch/v7"
    "github.com/ShopOnGO/search-service/configs"
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
        log.Fatalf("[elastic] failed to create client: %v", err)
    }

    // Проверка соединения
    res, err := client.Info()
    if err != nil {
        log.Fatalf("[elastic] failed to connect to Elasticsearch: %v", err)
    }
    defer res.Body.Close()

    ESClient = client
    Index = cfg.Elastic.Index

    log.Printf("[elastic] connected to %s | index=%s", cfg.Elastic.URL, Index)
}
