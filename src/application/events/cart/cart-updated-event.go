package events

import (
	"time"

	"cart/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartUpdatedEvent struct {
	ID        primitive.ObjectID `json:"id"`
	Products  []*models.Product  `json:"products"`
	Shipping  float32            `json:"shipping"`
	Discount  float32            `json:"discount"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Version   uint               `json:"version"`
}
