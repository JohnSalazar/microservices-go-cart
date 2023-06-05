package commands

import (
	"cart/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FinalizeCartCommand struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID `json:"customerId"`
	CouponID   primitive.ObjectID `json:"couponId"`
	Products   []*models.Product  `json:"products"`
	Shipping   float32            `json:"shipping"`
	Discount   float32            `json:"discount"`
	CardNumber string             `json:"cardNumber"`
	Version    uint               `json:"version"`
}
