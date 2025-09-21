package search

type ESProduct struct {
	ID          uint      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Material    string      `json:"material,omitempty"`
	IsActive    bool        `json:"is_active"`
	ImageURLs   []string    `json:"image_urls,omitempty"`
	VideoURLs   []string    `json:"video_urls,omitempty"`
	CategoryID  uint        `json:"category_id"`
	BrandID     uint        `json:"brand_id"`
	Variants    []ESVariant `json:"variants"`
}

type ESVariant struct {
	VariantID 		uint   	`json:"variant_id"`
	SKU       		string   	`json:"sku"`
	Price     		float64  	`json:"price"`
	Stock     		int      	`json:"stock"`
	Reserved_stock 	int 		`json:"reserved_stock"`
	Rating    		float64  	`json:"rating"`
	Discount    	float64  	`json:"discount,omitempty"`
	// FinalPrice  	float64  	`json:"final_price,omitempty"`
	ImageURLs   	[]string 	`json:"image_urls,omitempty"`
	Barcode     	string   	`json:"barcode,omitempty"`
	Dimensions  	string   	`json:"dimensions,omitempty"`
	MinOrder    	int      	`json:"min_order,omitempty"`
	IsActive    	bool     	`json:"is_active"`
	ReviewCount 	int      	`json:"review_count,omitempty"`
	Sizes       	string   	`json:"sizes"`
	Colors      	string   	`json:"colors"`
}

type SearchResponse struct {
	Products []ESProduct `json:"products"`
	Total    int         `json:"total"`
	Page     int         `json:"page"`
	Limit    int         `json:"limit"`
	Pages    int         `json:"pages"`
}