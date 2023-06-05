package dtos

import (
	"cart/src/models"
)

type FinalizeCart struct {
	ID         string            `json:"id"`
	CouponID   string            `json:"couponId"`
	Products   []*models.Product `json:"products"`
	Shipping   float32           `json:"shipping"`
	Discount   float32           `json:"discount"`
	CardNumber []byte            `json:"cardNumber"`
	Version    uint              `json:"version"`
}
