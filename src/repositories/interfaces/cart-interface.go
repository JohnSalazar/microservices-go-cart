package interfaces

import (
	"cart/src/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartRepository interface {
	FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) (*models.Cart, error)
	FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Cart, error)
	Create(ctx context.Context, cart *models.Cart) (*models.Cart, error)
	Update(ctx context.Context, cart *models.Cart) (*models.Cart, error)
	Delete(ctx context.Context, ID primitive.ObjectID) error
}
