package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ShopOnGO/search-service/internal/elastic"
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type SearchHandler struct{}

func NewSearchHandler(router *gin.Engine) *SearchHandler {
	handler := &SearchHandler{}

	searchGroup := router.Group("/search-service/products")
	{
		searchGroup.GET("/search", handler.searchProducts)
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
	name := c.Query("name")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	// Построим Elasticsearch Query DSL
	var mustClauses []map[string]interface{}

	if name != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				"name": name,
			},
		})
	}

	if minPrice != "" || maxPrice != "" {
		rangeQuery := map[string]interface{}{
			"range": map[string]interface{}{
				"variants.price": map[string]interface{}{},
			},
		}
		if minPrice != "" {
			if p, err := strconv.ParseFloat(minPrice, 64); err == nil {
				rangeQuery["range"].(map[string]interface{})["variants.price"].(map[string]interface{})["gte"] = p
			}
		}
		if maxPrice != "" {
			if p, err := strconv.ParseFloat(maxPrice, 64); err == nil {
				rangeQuery["range"].(map[string]interface{})["variants.price"].(map[string]interface{})["lte"] = p
			}
		}
		mustClauses = append(mustClauses, rangeQuery)
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации запроса"})
		return
	}

	res, err := elastic.ESClient.Search(
		elastic.ESClient.Search.WithIndex(elastic.Index),
		elastic.ESClient.Search.WithBody(&buf),
		elastic.ESClient.Search.WithTrackTotalHits(true),
		elastic.ESClient.Search.WithPretty(),
	)
	if err != nil || res.StatusCode != http.StatusOK {
		log.Printf("Search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка поиска"})
		return
	}
	defer res.Body.Close()

	var esResp struct {
		Hits struct {
			Hits []struct {
				Source ESProduct `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(bodyBytes, &esResp); err != nil {
		log.Printf("Unmarshal error: %v\nRaw: %s", err, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обработки ответа ES"})
		return
	}

	var products []ESProduct
	for _, hit := range esResp.Hits.Hits {
		products = append(products, hit.Source)
	}

	c.JSON(http.StatusOK, products)
}
