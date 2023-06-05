package dtos

type UpdateCoupon struct {
	ID           string  `json:"id"`
	Name         string  `bson:"name" json:"name"`
	Value        float32 `bson:"value" json:"value"`
	IsPercentage bool    `bson:"isPercentage" json:"isPercentage"`
	Quantity     uint    `bson:"quantity" json:"quantity"`
	Active       bool    `bson:"active" json:"active"`
	Version      uint    `json:"version"`
}
