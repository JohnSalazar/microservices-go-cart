package validators

import (
	"cart/src/dtos"

	"github.com/JohnSalazar/microservices-go-common/helpers"
	common_validator "github.com/JohnSalazar/microservices-go-common/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type addCoupon struct {
	Name     string  `from:"name" json:"name" validate:"required,max=150"`
	Value    float32 `from:"value" json:"value" validate:"required,gte=1"`
	Quantity uint    `from:"quantity" json:"quantity" validate:"required,gte=1"`
}

type updateCoupon struct {
	ID       primitive.ObjectID `from:"id" json:"id" validate:"required"`
	Name     string             `from:"name" json:"name" validate:"required,max=150"`
	Value    float32            `from:"value" json:"value" validate:"required,gte=1"`
	Quantity uint               `from:"quantity" json:"quantity" validate:"gte=0"`
}

func ValidateAddCoupon(fields *dtos.AddCoupon) interface{} {
	addCoupon := addCoupon{
		Name:     fields.Name,
		Value:    fields.Value,
		Quantity: fields.Quantity,
	}

	err := common_validator.Validate(addCoupon)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUpdateCoupon(fields *dtos.UpdateCoupon) interface{} {
	updateCoupon := updateCoupon{
		ID:       helpers.StringToID(fields.ID),
		Name:     fields.Name,
		Value:    fields.Value,
		Quantity: fields.Quantity,
	}

	err := common_validator.Validate(updateCoupon)
	if err != nil {
		return err
	}

	return nil
}
