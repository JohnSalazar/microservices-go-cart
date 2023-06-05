package dtos

import (
	"cart/src/models"
)

type AddCart struct {
	CustomerID string            `json:"customerId"`
	Products   []*models.Product `json:"products"`
	Shipping   float32           `json:"shipping"`
	Discount   float32           `json:"discount"`
}
