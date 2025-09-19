package graph

import (
	"strconv"

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

	// Category/Brand cast
	cat := int32(p.CategoryID)
	brand := int32(p.BrandID)

	// Variants
	var variants []*model.Variant
	for _, v := range p.Variants {
		// sizes convert []int -> []int32
		sizes32 := make([]int32, 0, len(v.Sizes))
		for _, s := range v.Sizes {
			sizes32 = append(sizes32, int32(s))
		}

		var matPtr *string
		if v.Material != "" {
			matPtr = &v.Material
		}

		variants = append(variants, &model.Variant{
			VariantID: v.VariantID,
			Sku:       v.SKU,
			Price:     v.Price,
			Sizes:     sizes32,
			Colors:    v.Colors,
			Material:  matPtr,
			Stock:     int32(v.Stock),
			Rating:    v.Rating,
		})
	}

	// ID в model.Product — string
	idStr := strconv.FormatUint(uint64(p.ID), 10)

	return &model.Product{
		ID:          idStr,
		Name:        p.Name,
		Description: descPtr,
		CategoryID:  cat,
		BrandID:     brand,
		Variants:    variants,
	}
}
