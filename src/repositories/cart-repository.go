package repositories

import (
	"context"
	"encoding/json"
	"time"

	"cart/src/models"

	"github.com/google/uuid"
	"github.com/oceano-dev/microservices-go-common/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cartRepository struct {
	database *mongo.Database
}

func NewCartRepository(
	database *mongo.Database,
) *cartRepository {
	return &cartRepository{
		database: database,
	}
}

func (r *cartRepository) collectionName() string {
	return "carts"
}

func (r *cartRepository) collection() *mongo.Collection {
	return r.database.Collection(r.collectionName())
}

func (r *cartRepository) findOne(ctx context.Context, filter interface{}) (*models.Cart, error) {
	findOneOptions := options.FindOneOptions{}
	findOneOptions.SetSort(bson.M{"version": -1})

	newFilter := map[string]interface{}{
		"deleted": false,
	}
	mergeFilter := helpers.MergeFilters(newFilter, filter)

	object := map[string]interface{}{}
	err := r.collection().FindOne(ctx, mergeFilter, &findOneOptions).Decode(object)
	if err != nil {
		return nil, err
	}

	cart, err := r.mapCart(object)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *cartRepository) findOneAndUpdate(ctx context.Context, filter interface{}, fields interface{}) *mongo.SingleResult {
	findOneAndUpdateOptions := options.FindOneAndUpdateOptions{}
	findOneAndUpdateOptions.SetReturnDocument(options.After)

	result := r.collection().FindOneAndUpdate(ctx, filter, bson.M{"$set": fields}, &findOneAndUpdateOptions)

	return result
}

func (r *cartRepository) FindByCustomerID(ctx context.Context, customerID primitive.ObjectID) (*models.Cart, error) {
	filter := bson.M{"customer_id": customerID}

	return r.findOne(ctx, filter)
}

func (r *cartRepository) FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Cart, error) {
	filter := bson.M{"_id": ID}

	return r.findOne(ctx, filter)
}

func (r *cartRepository) Create(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
	products := r.mapCartProducts(cart.Products)

	fields := bson.M{
		"_id":         cart.ID,
		"customer_id": cart.CustomerID,
		"products":    products,
		"shipping":    cart.Shipping,
		"discount":    cart.Discount,
		"created_at":  time.Now().UTC(),
		"version":     0,
		"deleted":     false,
	}

	_, err := r.collection().InsertOne(ctx, fields)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *cartRepository) Update(ctx context.Context, cart *models.Cart) (*models.Cart, error) {
	cart.Version++
	cart.UpdatedAt = time.Now().UTC()

	products := r.mapCartProducts(cart.Products)

	fields := bson.M{
		"products":   products,
		"shipping":   cart.Shipping,
		"discount":   cart.Discount,
		"updated_at": cart.UpdatedAt,
		"version":    cart.Version,
	}

	filter := r.filterUpdate(cart)

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return nil, result.Err()
	}

	object := map[string]interface{}{}
	err := result.Decode(object)
	if err != nil {
		return nil, err
	}

	modelCart, err := r.mapCart(object)
	if err != nil {
		return nil, err
	}

	return modelCart, err
}

func (r *cartRepository) Delete(ctx context.Context, ID primitive.ObjectID) error {
	filter := bson.M{"_id": ID}

	fields := bson.M{"deleted": true}

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *cartRepository) filterUpdate(cart *models.Cart) interface{} {
	filter := bson.M{
		"_id":     cart.ID,
		"version": cart.Version - 1,
	}

	return filter
}

func (r *cartRepository) mapCart(object map[string]interface{}) (*models.Cart, error) {
	jsonStr, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	var cart models.Cart
	if err := json.Unmarshal(jsonStr, &cart); err != nil {
		return nil, err
	}

	cart.ID = object["_id"].(primitive.ObjectID)
	cart.CustomerID = object["customer_id"].(primitive.ObjectID)

	if object["products"] != nil {
		var products []*models.Product
		listProducts := object["products"].(primitive.A)
		for _, product := range listProducts {
			product, err := r.mapProductFromInterfaceToModel(product.(map[string]interface{}))
			if err != nil {
				return nil, err
			}

			products = append(products, product)
		}

		cart.Products = products
	}

	return &cart, nil
}

func (r *cartRepository) mapProductFromInterfaceToModel(object map[string]interface{}) (*models.Product, error) {
	jsonStr, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := json.Unmarshal(jsonStr, &product); err != nil {
		return nil, err
	}

	ID, err := uuid.Parse(object["_id"].(string))
	if err != nil {
		return nil, err
	}
	product.ID = ID

	return &product, nil
}

func (r *cartRepository) mapCartProducts(cartProducts []*models.Product) []map[string]interface{} {
	var products []map[string]interface{}
	for _, product := range cartProducts {
		modelProduct := map[string]interface{}{
			"_id":         product.ID.String(),
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"quantity":    product.Quantity,
			"image":       product.Image,
		}

		products = append(products, modelProduct)
	}

	return products
}
