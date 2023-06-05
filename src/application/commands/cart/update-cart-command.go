package commands

import (
	"cart/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateCartCommand struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID `json:"customerId"`
	Products   []*models.Product  `json:"products"`
	Shipping   float32            `json:"shipping"`
	Discount   float32            `json:"discount"`
	Version    uint               `json:"version"`
}
