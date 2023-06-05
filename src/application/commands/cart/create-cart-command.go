package commands

import (
	"time"

	"cart/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCartCommand struct {
	ID         primitive.ObjectID `json:"id"`
	CustomerID primitive.ObjectID `json:"customerId"`
	Products   []*models.Product  `json:"products"`
	Shipping   float32            `json:"shipping"`
	Discount   float32            `json:"discount"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	Version    uint               `json:"version"`
	Deleted    bool               `json:"deleted,omitempty"`
}
