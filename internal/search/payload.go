package search

type ProductCreatedEvent struct {
	Action      string                 `json:"action"`
	ProductID   uint                   `json:"product_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Material    string                 `json:"material"`
	IsActive    bool                   `json:"is_active"`
	CategoryID  uint                   `json:"category_id"`
	BrandID     uint                   `json:"brand_id"`

	Variants    []ProductVariantInput  `json:"variants"`
}

type ProductVariantInput struct {
	VariantID      uint     `json:"variant_id"`
	SKU            string   `json:"sku"`
	Price          float64  `json:"price"`
	Discount       float64  `json:"discount"`
	Stock          int      `json:"stock"`
	Reserved_stock int      `json:"reserved_stock"`
	Rating         float64  `json:"rating"`
	ImageURLs      []string `json:"images"` 
	Sizes          string   `json:"sizes"`
	Colors         string   `json:"colors"`
	Barcode        string   `json:"barcode"`
	Dimensions     string   `json:"dimensions"`
	MinOrder       int      `json:"min_order"`
	IsActive       bool     `json:"is_active"`
	ReviewCount    int      `json:"review_count"`
}