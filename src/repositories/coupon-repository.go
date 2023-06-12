package repositories

import (
	"context"
	"time"

	"cart/src/models"

	"github.com/JohnSalazar/microservices-go-common/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type couponRepository struct {
	database *mongo.Database
}

func NewCouponRepository(
	database *mongo.Database,
) *couponRepository {
	return &couponRepository{
		database: database,
	}
}

func (r *couponRepository) collectionName() string {
	return "coupons"
}

func (r *couponRepository) collection() *mongo.Collection {
	return r.database.Collection(r.collectionName())
}

func (r *couponRepository) find(ctx context.Context, filter interface{}, page int, size int) ([]*models.Coupon, error) {
	findOptions := options.FindOptions{}
	findOptions.SetSort(bson.M{"name": 1})

	page64 := int64(page)
	size64 := int64(size)
	findOptions.SetSkip((page64 - 1) * size64)
	findOptions.SetLimit(size64)

	newFilter := map[string]interface{}{
		"deleted": false,
	}
	mergeFilter := helpers.MergeFilters(newFilter, filter)

	cursor, err := r.collection().Find(ctx, mergeFilter, &findOptions)
	if err != nil {
		defer cursor.Close(ctx)
		return nil, err
	}

	var coupons []*models.Coupon

	for cursor.Next(ctx) {
		coupon := &models.Coupon{}

		err = cursor.Decode(coupon)
		if err != nil {
			return nil, err
		}

		coupons = append(coupons, coupon)
	}

	return coupons, nil
}

func (r *couponRepository) findOne(ctx context.Context, filter interface{}) (*models.Coupon, error) {
	findOneOptions := options.FindOneOptions{}
	findOneOptions.SetSort(bson.M{"version": -1})

	newFilter := map[string]interface{}{
		"deleted": false,
	}
	mergeFilter := helpers.MergeFilters(newFilter, filter)

	coupon := &models.Coupon{}
	err := r.collection().FindOne(ctx, mergeFilter, &findOneOptions).Decode(coupon)
	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (r *couponRepository) findOneAndUpdate(ctx context.Context, filter interface{}, fields interface{}) *mongo.SingleResult {
	findOneAndUpdateOptions := options.FindOneAndUpdateOptions{}
	findOneAndUpdateOptions.SetReturnDocument(options.After)

	result := r.collection().FindOneAndUpdate(ctx, filter, bson.M{"$set": fields}, &findOneAndUpdateOptions)

	return result
}

func (r *couponRepository) GetAll(ctx context.Context, name string, page int, size int) ([]*models.Coupon, error) {
	filter := bson.M{
		"name":     bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}},
		"quantity": bson.M{"$gt": 0}}

	return r.find(ctx, filter, page, size)
}

func (r *couponRepository) GetByName(ctx context.Context, name string) (*models.Coupon, error) {
	filter := bson.M{"name": name, "active": true, "quantity": bson.M{"$gt": 0}}

	return r.findOne(ctx, filter)
}

func (r *couponRepository) FindByName(ctx context.Context, name string) (*models.Coupon, error) {
	filter := bson.M{"name": name}

	return r.findOne(ctx, filter)
}

func (r *couponRepository) FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Coupon, error) {
	filter := bson.M{"_id": ID}

	return r.findOne(ctx, filter)
}

func (r *couponRepository) Create(ctx context.Context, coupon *models.Coupon) (*models.Coupon, error) {
	fields := bson.M{
		"_id":          coupon.ID,
		"name":         coupon.Name,
		"value":        coupon.Value,
		"isPercentage": coupon.IsPercentage,
		"quantity":     coupon.Quantity,
		"active":       coupon.Active,
		"created_at":   time.Now().UTC(),
		"version":      0,
		"deleted":      false,
	}

	_, err := r.collection().InsertOne(ctx, fields)

	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func (r *couponRepository) Update(ctx context.Context, coupon *models.Coupon) (*models.Coupon, error) {
	coupon.Version++
	coupon.UpdatedAt = time.Now().UTC()

	fields := bson.M{
		"name":         coupon.Name,
		"value":        coupon.Value,
		"isPercentage": coupon.IsPercentage,
		"quantity":     coupon.Quantity,
		"active":       coupon.Active,
		"updated_at":   coupon.UpdatedAt,
		"version":      coupon.Version,
	}

	filter := r.filterUpdate(coupon)

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return nil, result.Err()
	}

	modelCoupon := &models.Coupon{}
	decodeErr := result.Decode(modelCoupon)

	return modelCoupon, decodeErr
}

func (r *couponRepository) Delete(ctx context.Context, ID primitive.ObjectID) error {
	filter := bson.M{"_id": ID}

	fields := bson.M{"deleted": true}

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *couponRepository) filterUpdate(coupon *models.Coupon) interface{} {
	filter := bson.M{
		"_id":     coupon.ID,
		"version": coupon.Version - 1,
	}

	return filter
}
