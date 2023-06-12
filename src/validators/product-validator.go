package validators

import (
	"cart/src/dtos"

	common_validator "github.com/JohnSalazar/microservices-go-common/validators"
	"github.com/google/uuid"
)

type product struct {
	ID       uuid.UUID `from:"id" json:"id" validate:"required"`
	Name     string    `from:"name" json:"name" validate:"required,max=150"`
	Price    float32   `from:"price" json:"price" validate:"required,gte=1"`
	Quantity uint      `from:"quantity" json:"quantity" validate:"required,gte=1"`
}

func ValidateProduct(fields *dtos.Product) interface{} {
	product := product{
		ID:       fields.ID,
		Name:     fields.Name,
		Price:    fields.Price,
		Quantity: fields.Quantity,
	}

	err := common_validator.Validate(product)
	if err != nil {
		return err
	}

	return nil
}
