package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cart struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customerId,omitempty"`
	Products   []*Product         `bson:"products" json:"products"`
	Shipping   float32            `bson:"shipping" json:"shipping"`
	Discount   float32            `bson:"discount" json:"discount"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	Version    uint               `bson:"version" json:"version"`
	Deleted    bool               `bson:"deleted" json:"deleted,omitempty"`
}

func (c *Cart) SumProduct() float32 {
	var sum float32 = 0
	for _, product := range c.Products {
		sum += product.Sum()
	}

	return sum
}

func (c *Cart) ApplyCoupon(coupon *Coupon) float32 {
	var discount float32 = 0
	if coupon.IsPercentage {
		if coupon.Value > 100 {
			coupon.Value = 100
		}
	}

	if coupon.Value > 0 {
		sum := c.SumProduct()
		if coupon.IsPercentage {
			discount = (sum * coupon.Value) / 100
		} else {
			discount = sum - coupon.Value
			if coupon.Value > sum {
				discount = sum
			}
		}
	}

	c.Discount = discount

	return discount
}

func (c *Cart) Total() float32 {
	return c.SumProduct() - c.Discount + c.Shipping
}
