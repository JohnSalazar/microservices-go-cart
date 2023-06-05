package events

import (
	"time"

	"cart/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartCreatedEvent struct {
	ID        primitive.ObjectID `json:"id"`
	Products  []*models.Product  `json:"products"`
	Shipping  float32            `json:"shipping"`
	Discount  float32            `json:"discount"`
	CreatedAt time.Time          `json:"createdAt"`
	Version   uint               `json:"version"`
}
