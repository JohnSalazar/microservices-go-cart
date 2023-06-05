package commands

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdateCouponCommand struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Value        float32            `json:"value"`
	IsPercentage bool               `json:"isPercentage"`
	Quantity     uint               `json:"quantity"`
	Active       bool               `json:"active"`
	Version      uint               `json:"version"`
}
