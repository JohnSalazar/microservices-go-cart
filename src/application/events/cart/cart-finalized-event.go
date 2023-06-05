package events

import (
	"cart/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartFinalizedEvent struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID `json:"customerId"`
	Products   []*models.Product  `json:"products"`
	Sum        float32            `json:"sum"`
	Shipping   float32            `json:"shipping"`
	Discount   float32            `json:"discount"`
	CardNumber []byte             `json:"cardNumber"`
	Kid        string             `json:"kid"`
}
