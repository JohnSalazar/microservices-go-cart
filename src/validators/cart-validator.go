package validators

import (
	"cart/src/dtos"
	"cart/src/models"

	"github.com/JohnSalazar/microservices-go-common/helpers"
	common_validator "github.com/JohnSalazar/microservices-go-common/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type addCart struct {
	CustomerID primitive.ObjectID `from:"customerId" json:"customerId" validate:"required"`
	Products   []*models.Product  `from:"products" json:"products" validate:"required"`
}

type updateCart struct {
	ID       primitive.ObjectID `from:"id" json:"id" validate:"required"`
	Products []*models.Product  `from:"products" json:"products" validate:"required"`
}

type finalizeCart struct {
	ID         primitive.ObjectID `from:"id" json:"id" validate:"required"`
	Products   []*models.Product  `from:"products" json:"products" validate:"required"`
	CardNumber []byte             `from:"cardNumber" json:"cardNumber" validate:"required"`
}

type deleteCart struct {
	ID primitive.ObjectID `from:"id" json:"id" validate:"required"`
}

func ValidateAddCart(fields *dtos.AddCart) interface{} {
	addCart := addCart{
		CustomerID: helpers.StringToID(fields.CustomerID),
		Products:   fields.Products,
	}

	err := common_validator.Validate(addCart)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUpdateCart(fields *dtos.UpdateCart) interface{} {
	updateCart := updateCart{
		ID:       helpers.StringToID(fields.ID),
		Products: fields.Products,
	}

	err := common_validator.Validate(updateCart)
	if err != nil {
		return err
	}

	return nil
}

func ValidateFinalizeCart(fields *dtos.FinalizeCart) interface{} {
	finalizeCart := finalizeCart{
		ID:         helpers.StringToID(fields.ID),
		Products:   fields.Products,
		CardNumber: fields.CardNumber,
	}

	err := common_validator.Validate(finalizeCart)
	if err != nil {
		return err
	}

	return nil
}

func ValidateDeleteCart(fields *dtos.DeleteCart) interface{} {
	deleteCart := deleteCart{
		ID: helpers.StringToID(fields.ID),
	}

	err := common_validator.Validate(deleteCart)
	if err != nil {
		return err
	}

	return nil
}
