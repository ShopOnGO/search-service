package product

import (
	"errors"
	"strings"
)

type SearchRequest struct {
	Name     string  `form:"name" binding:"omitempty,min=1,max=100"`
	MinPrice float64 `form:"min_price" binding:"omitempty,min=0, regexp=^\\d+(\\.\\d+)?$"`
	MaxPrice float64 `form:"max_price" binding:"omitempty,min=0, regexp=^\\d+(\\.\\d+)?$"`
	Page     int     `form:"page" binding:"omitempty,min=1, regexp=^\\d+$"`
	Limit    int     `form:"limit" binding:"omitempty,min=1,max=100, regexp=^\\d+$"`
}

func (r *SearchRequest) Validate() error {
	if r.MinPrice > 0 && r.MaxPrice > 0 && r.MinPrice > r.MaxPrice {
		return errors.New("min_price cannot be greater than max_price")
	}

	if r.MaxPrice > 0 && r.MaxPrice < r.MinPrice {
		return errors.New("max_price cannot be less than min_price")
	}

	if r.Name != "" {
		r.Name = strings.TrimSpace(r.Name)
		if len(r.Name) == 0 {
			return errors.New("name cannot be empty")
		}
	}

	if r.Page < 0 {
		return errors.New("page must be positive")
	}

	if r.Limit < 0 {
		return errors.New("limit must be positive")
	}

	return nil
}
