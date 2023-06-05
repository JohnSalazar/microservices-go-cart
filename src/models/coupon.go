package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coupon struct {
	ID           primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Value        float32            `bson:"value" json:"value"`
	IsPercentage bool               `bson:"isPercentage" json:"isPercentage"`
	Quantity     uint               `bson:"quantity" json:"quantity"`
	Active       bool               `bson:"active" json:"active"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	Version      uint               `bson:"version" json:"version"`
	Deleted      bool               `bson:"deleted" json:"deleted,omitempty"`
}
