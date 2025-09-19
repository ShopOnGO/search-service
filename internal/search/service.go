package search

import "github.com/ShopOnGO/ShopOnGO/pkg/logger"

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) AddProductToElasticSearch(product ProductCreatedEvent) error {
	logger.Infof("Добавление продукта в индекс Elasticsearch: %+v", product)
	return nil
}
