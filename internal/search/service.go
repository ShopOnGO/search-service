package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/search-service/internal/elastic" // ← где лежат ESClient и Index
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) AddProductToElasticSearch(product ProductCreatedEvent) error {
	logger.Infof("Добавление продукта в индекс Elasticsearch: %+v", product)

	// Проверка: ProductID обязателен
	if product.ProductID == 0 {
		return fmt.Errorf("product_id is required")
	}

	// Формируем документ для Elasticsearch
	doc := map[string]interface{}{
		"id":          fmt.Sprintf("%d", product.ProductID), // ← keyword в маппинге
		"name":        product.Name,
		"description": product.Description,
		"category_id": product.CategoryID,
		"brand_id":    product.BrandID,
		"image_urls":  product.ImageKeys,
		"video_urls":  product.VideoKeys,

		// Варианты — создаём хотя бы один по умолчанию
		"variants": []map[string]interface{}{
			{
				"variant_id":   fmt.Sprintf("variant-%d-default", product.ProductID),
				"sku":          fmt.Sprintf("SKU-%d", product.ProductID),
				"price":        float64(product.Price),
				"discount":     float64(product.Discount),
				"stock":        999,
				"min_order":    1,
				"sizes":        []int{},
				"colors":       []string{},
				"material":     "",
				"rating":       0.0,
				"review_count": 0,
				"rating_sum":   0,
				"dimensions":   "",
			},
		},
	}

	// Сериализуем документ в JSON
	body, err := json.Marshal(doc)
	if err != nil {
		logger.Errorf("Ошибка сериализации документа: %v", err)
		return fmt.Errorf("failed to serialize document: %w", err)
	}

	// Индексируем в Elasticsearch
	req := esapi.IndexRequest{
		Index:      elastic.Index,
		DocumentID: fmt.Sprintf("%d", product.ProductID), // ← _id документа
		Body:       bytes.NewReader(body),
		Refresh:    "true", // ← опционально: сразу видно в поиске (для dev) не надо Refresh: "true" в продакшене
	}

	res, err := req.Do(context.Background(), elastic.ESClient)
	if err != nil {
		logger.Errorf("Ошибка запроса к Elasticsearch: %v", err)
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			logger.Errorf("Ошибка парсинга ошибки Elasticsearch: %v", err)
		} else {
			logger.Errorf("Ошибка Elasticsearch: %s", e)
		}
		return fmt.Errorf("Elasticsearch error: %s", res.Status())
	}

	logger.Infof("✅ Продукт ID %d успешно проиндексирован в Elasticsearch", product.ProductID)
	return nil
}
