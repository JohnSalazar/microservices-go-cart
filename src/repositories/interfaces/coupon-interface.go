package interfaces

import (
	"cart/src/models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouponRepository interface {
	GetAll(ctx context.Context, name string, page int, size int) ([]*models.Coupon, error)
	GetByName(ctx context.Context, name string) (*models.Coupon, error)
	FindByName(ctx context.Context, name string) (*models.Coupon, error)
	FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Coupon, error)
	Create(ctx context.Context, coupon *models.Coupon) (*models.Coupon, error)
	Update(ctx context.Context, coupon *models.Coupon) (*models.Coupon, error)
	Delete(ctx context.Context, ID primitive.ObjectID) error
}
