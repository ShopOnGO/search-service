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
		"id":          product.ProductID,
		"name":        product.Name,
		"description": product.Description,
		"category_id": product.CategoryID,
		"brand_id":    product.BrandID,
		"material":    product.Material,
    	"is_active":   product.IsActive,
	}

	var variants []map[string]interface{}
    for _, v := range product.Variants {
        variants = append(variants, map[string]interface{}{
            "variant_id":   v.VariantID,
            "sku":          v.SKU,
            "price":        v.Price,
            "discount":     v.Discount,
            // "final_price":  v.Price - v.Discount,
            "sizes":        v.Sizes,
            "colors":       v.Colors,
            "stock":        v.Stock,
            "rating":       v.Rating,
            "review_count": v.ReviewCount,
            "barcode":      v.Barcode,
            "dimensions":   v.Dimensions,
            "image_urls":   v.ImageURLs,
            "min_order":    v.MinOrder,
            "is_active":    v.IsActive,
        })
    }
    doc["variants"] = variants

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
		Refresh:    "true", // для dev — в продакшене убери или сделай опционально
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
