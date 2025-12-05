package graph

import (
	"github.com/ShopOnGO/search-service/internal/graph/model"
	"github.com/ShopOnGO/search-service/internal/search"
)

// Resolver root for gqlgen.
type Resolver struct{}

func ConvertESProductToModel(p *search.ESProduct) *model.Product {
	// Description pointer
	var descPtr *string
	if p.Description != "" {
		descPtr = &p.Description
	}

	// Variants
	var variants []*model.Variant
	for _, v := range p.Variants {
		variants = append(variants, &model.Variant{
			VariantID:   int32(v.VariantID),
			Sku:         v.SKU,
			Price:       v.Price,
			Discount:    v.Discount,
			// FinalPrice:  v.Price - v.Discount,
			Sizes:       v.Sizes,
			Colors:      v.Colors,
			Stock:       int32(v.Stock),
			Rating:      v.Rating,
			ReviewCount: int32(v.ReviewCount),
			Barcode:     &v.Barcode,
			Dimensions:  &v.Dimensions,
			MinOrder:    int32(v.MinOrder),
			IsActive:    v.IsActive,
			ImageUrls:   v.ImageURLs, 
		})
	}


	return &model.Product{
		ID:          int32(p.ID),
		Name:        p.Name,
		Description: descPtr,
		Material:    &p.Material,
		CategoryID:  int32(p.CategoryID),
		BrandID:     int32(p.BrandID),
		IsActive:    p.IsActive,
		ImageUrls:   p.ImageURLs,
		VideoUrls:   p.VideoURLs,
		Variants:    variants,
	}
}
