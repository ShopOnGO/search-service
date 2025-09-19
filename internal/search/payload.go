package search

type ProductCreatedEvent struct {
	Action    	string 		`json:"action"`
	ProductID 	uint   		`json:"product_id"`

	Name        string  	`json:"name"`
	Description string  	`json:"description"`
	Price       int64   	`json:"price"`
	Discount    int64   	`json:"discount"`

	CategoryID  uint    	`json:"category_id"`
	BrandID     uint    	`json:"brand_id"`

	ImageKeys  	[]string	`json:"image_keys"`
	VideoKeys  	[]string 	`json:"video_keys"`
}