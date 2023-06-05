package dtos

import (
	"cart/src/models"
)

type UpdateCart struct {
	ID       string            `json:"id"`
	Products []*models.Product `json:"products"`
	Shipping float32           `json:"shipping"`
	Discount float32           `json:"discount"`
	Version  uint              `json:"version"`
}
