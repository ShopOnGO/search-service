package search

type ESVariant struct {
	VariantID string   `json:"variant_id"`
	SKU       string   `json:"sku"`
	Price     float64  `json:"price"`
	Sizes     []int    `json:"sizes"`
	Colors    []string `json:"colors"`
	Material  string   `json:"material"`
	Stock     int      `json:"stock"`
	Rating    float64  `json:"rating"`
}

type ESProduct struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	CategoryID  uint        `json:"category_id"`
	BrandID     uint        `json:"brand_id"`
	Variants    []ESVariant `json:"variants"`
}

type SearchResponse struct {
	Products []ESProduct `json:"products"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	Limit    int         `json:"limit"`
	Pages    int         `json:"pages"`
}