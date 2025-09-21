package search

type ProductCreatedEvent struct {
	Action    	string 		`json:"action"`
	ProductID 	uint   		`json:"product_id"`

	Name        string  	`json:"name"`
	Description string  	`json:"description"`
	Material    string      `json:"material"`
	IsActive    bool        `json:"is_active"`
	ImageURLs   []string    `json:"image_urls"`
	VideoURLs   []string    `json:"video_urls"`

	CategoryID  uint    	`json:"category_id"`
	BrandID     uint    	`json:"brand_id"`
	
	ImageKeys  	[]string	`json:"image_keys"`
	VideoKeys  	[]string 	`json:"video_keys"`

	Variants    []ESVariant `json:"variants"`
}