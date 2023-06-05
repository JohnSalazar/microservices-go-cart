package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `bson:"_id" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	Price       float32   `bson:"price" json:"price"`
	Quantity    uint      `bson:"quantity" json:"quantity"`
	Image       string    `bson:"image" json:"image"`
}

func (p Product) Sum() float32 {
	// price := p.Price
	// quantity := float32(p.Quantity)
	// return price * quantity
	return p.Price * float32(p.Quantity)
}
