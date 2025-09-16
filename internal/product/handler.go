package product

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/ShopOnGO/search-service/internal/elastic"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct{}

func NewSearchHandler(router *gin.Engine) *SearchHandler {
	handler := &SearchHandler{}

	searchGroup := router.Group("/search-service/products")
	{
		searchGroup.GET("/search", handler.searchProducts)
		searchGroup.GET("/", handler.getAllProducts)
	}

	return handler
}

// searchProducts godoc
// @Summary Поиск продуктов
// @Description Ищет продукты по ключевому слову (name) и фильтрам
// @Tags Поиск
// @Param name query string false "Название продукта"
// @Param min_price query number false "Минимальная цена"
// @Param max_price query number false "Максимальная цена"
// @Success 200 {array} product.ESProduct
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /search-service/products/search [get]
func (h *SearchHandler) searchProducts(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters", "details": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	// Построим Elasticsearch Query DSL
	var mustClauses []map[string]interface{}

	// Поиск по названию (поддержка частичного совпадения)
	if req.Name != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     req.Name,
				"fields":    []string{"name^2", "description"},
				"type":      "best_fields",
				"fuzziness": "AUTO",
			},
		})
	}

	// Фильтр по цене
	if req.MinPrice > 0 || req.MaxPrice > 0 {
		rangeQuery := map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "variants",
				"query": map[string]interface{}{
					"range": map[string]interface{}{
						"variants.price": map[string]interface{}{},
					},
				},
			},
		}

		priceRange := make(map[string]interface{})
		if req.MinPrice > 0 {
			priceRange["gte"] = req.MinPrice
		}
		if req.MaxPrice > 0 {
			priceRange["lte"] = req.MaxPrice
		}

		rangeQuery["nested"].(map[string]interface{})["query"].(map[string]interface{})["range"].(map[string]interface{})["variants.price"] = priceRange
		mustClauses = append(mustClauses, rangeQuery)
	}

	// Если нет условий поиска, возвращаем все товары
	var query map[string]interface{}
	if len(mustClauses) == 0 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": mustClauses,
				},
			},
		}
	}

	// Добавляем пагинацию и сортировку
	query["from"] = (req.Page - 1) * req.Limit
	query["size"] = req.Limit
	query["sort"] = []map[string]interface{}{
		{"_score": map[string]interface{}{"order": "desc"}},
		{"name.keyword": map[string]interface{}{"order": "asc"}},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logger.Infof("Query serialization error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации запроса"})
		return
	}

	res, err := elastic.ESClient.Search(
		elastic.ESClient.Search.WithIndex(elastic.Index),
		elastic.ESClient.Search.WithBody(&buf),
		elastic.ESClient.Search.WithTrackTotalHits(true),
		elastic.ESClient.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Elasticsearch search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка подключения к поисковой системе"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Elasticsearch returned status %d: %s", res.StatusCode, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка поиска"})
		return
	}

	var esResp struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source ESProduct `json:"_source"`
				Score  float64   `json:"_score"`
			} `json:"hits"`
		} `json:"hits"`
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения ответа"})
		return
	}

	if err := json.Unmarshal(bodyBytes, &esResp); err != nil {
		log.Printf("Unmarshal error: %v\nRaw: %s", err, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обработки ответа поисковой системы"})
		return
	}

	// Формируем ответ с метаданными
	var products []ESProduct
	for _, hit := range esResp.Hits.Hits {
		products = append(products, hit.Source)
	}

	response := SearchResponse{
		Products: products,
		Total:    esResp.Hits.Total.Value,
		Page:     req.Page,
		Limit:    req.Limit,
		Pages:    (esResp.Hits.Total.Value + req.Limit - 1) / req.Limit,
	}

	c.JSON(http.StatusOK, response)
}

// getAllProducts godoc
// @Summary Получить все продукты
// @Description Возвращает все продукты с пагинацией (ограничено 100 элементами на страницу)
// @Tags Поиск
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(20) minimum(1) maximum(100)
// @Success 200 {object} product.SearchResponse
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /search-service/products/ [get]
func (h *SearchHandler) getAllProducts(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameters", "details": err.Error()})
		return
	}

	// Устанавливаем значения по умолчанию
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	if req.Limit > 100 { // Ограничение для защиты ES
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit cannot exceed 100 items per page",
			"max_limit": 100,
			"provided_limit": req.Limit,
		})
		return
	}

	// Простой запрос для получения всех товаров
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"from": (req.Page - 1) * req.Limit,
		"size": req.Limit,
		"sort": []map[string]interface{}{
			{"name.keyword": map[string]interface{}{"order": "asc"}},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Printf("Query serialization error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации запроса"})
		return
	}

	res, err := elastic.ESClient.Search(
		elastic.ESClient.Search.WithIndex(elastic.Index),
		elastic.ESClient.Search.WithBody(&buf),
		elastic.ESClient.Search.WithTrackTotalHits(true),
		elastic.ESClient.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Elasticsearch search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка подключения к поисковой системе"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Elasticsearch returned status %d: %s", res.StatusCode, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка поиска"})
		return
	}

	var esResp struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source ESProduct `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка чтения ответа"})
		return
	}

	if err := json.Unmarshal(bodyBytes, &esResp); err != nil {
		log.Printf("Unmarshal error: %v\nRaw: %s", err, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обработки ответа поисковой системы"})
		return
	}

	// Формируем ответ с метаданными
	var products []ESProduct
	for _, hit := range esResp.Hits.Hits {
		products = append(products, hit.Source)
	}

	response := SearchResponse{
		Products: products,
		Total:    esResp.Hits.Total.Value,
		Page:     req.Page,
		Limit:    req.Limit,
		Pages:    (esResp.Hits.Total.Value + req.Limit - 1) / req.Limit,
	}

	c.JSON(http.StatusOK, response)
}
