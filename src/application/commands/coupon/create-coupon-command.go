package commands

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateCouponCommand struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Value        float32            `json:"value"`
	IsPercentage bool               `json:"isPercentage"`
	Quantity     uint               `json:"quantity"`
	Active       bool               `json:"active"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty"`
	Version      uint               `json:"version"`
	Deleted      bool               `json:"deleted,omitempty"`
}
