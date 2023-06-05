package commands

import (
	"context"
	"errors"
	"strings"
	"time"

	events "cart/src/application/events/cart"
	"cart/src/dtos"
	"cart/src/models"
	"cart/src/repositories/interfaces"
	"cart/src/validators"

	common_security "github.com/oceano-dev/microservices-go-common/security"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartCommandHandler struct {
	cartRepository         interfaces.CartRepository
	cartEventHandler       *events.CartEventHandler
	managerSecurityRSAKeys common_security.ManagerSecurityRSAKeys
}

func NewCartCommandHandler(
	cartRepository interfaces.CartRepository,
	cartEventHandler *events.CartEventHandler,
	managerSecurityRSAKeys common_security.ManagerSecurityRSAKeys,
) *CartCommandHandler {
	return &CartCommandHandler{
		cartRepository:         cartRepository,
		cartEventHandler:       cartEventHandler,
		managerSecurityRSAKeys: managerSecurityRSAKeys,
	}
}

func (cart *CartCommandHandler) CreateCartCommandHandler(ctx context.Context, command *CreateCartCommand) (*models.Cart, error) {
	cartDto := &dtos.AddCart{
		CustomerID: command.CustomerID.Hex(),
		Products:   command.Products,
		Shipping:   command.Shipping,
		Discount:   command.Discount,
	}

	result := validators.ValidateAddCart(cartDto)
	if result != nil {
		return nil, errors.New(strings.Join(result.([]string), ""))
	}

	err := cart.validateProduct(command.Products)
	if err != nil {
		return nil, err
	}

	cartModel := &models.Cart{
		ID:         primitive.NewObjectID(),
		CustomerID: command.CustomerID,
		Products:   cartDto.Products,
		Shipping:   cartDto.Shipping,
		Discount:   cartDto.Discount,
		CreatedAt:  time.Now().UTC(),
	}

	cartExists, _ := cart.cartRepository.FindByCustomerID(ctx, cartModel.CustomerID)
	if cartExists != nil {
		return nil, errors.New("already a cart for this customer")
	}

	cartModel, err = cart.cartRepository.Create(ctx, cartModel)
	if err != nil {
		return nil, err
	}

	cartEvent := &events.CartCreatedEvent{
		ID:        cartModel.ID,
		Products:  cartModel.Products,
		Shipping:  cartModel.Shipping,
		Discount:  cartModel.Discount,
		CreatedAt: cartModel.CreatedAt,
		Version:   cartModel.Version,
	}

	go cart.cartEventHandler.CartCreatedEventHandler(ctx, cartEvent)

	return cartModel, nil
}

func (cart *CartCommandHandler) UpdateCartCommandHandler(ctx context.Context, command *UpdateCartCommand) (*models.Cart, error) {
	cartDto := &dtos.UpdateCart{
		ID:       command.ID.Hex(),
		Products: command.Products,
		Shipping: command.Shipping,
		Discount: command.Discount,
		Version:  command.Version,
	}

	result := validators.ValidateUpdateCart(cartDto)
	if result != nil {
		return nil, errors.New(strings.Join(result.([]string), ""))
	}

	err := cart.validateProduct(command.Products)
	if err != nil {
		return nil, err
	}

	cartModel := &models.Cart{
		ID:        command.ID,
		Products:  cartDto.Products,
		Shipping:  cartDto.Shipping,
		Discount:  cartDto.Discount,
		Version:   cartDto.Version,
		UpdatedAt: time.Now().UTC(),
	}

	cartExists, err := cart.cartRepository.FindByID(ctx, cartModel.ID)
	if err != nil {
		return nil, err
	}

	if command.CustomerID != cartExists.CustomerID {
		return nil, errors.New("you can only change your cart")
	}

	cartModel, err = cart.cartRepository.Update(ctx, cartModel)
	if err != nil {
		return nil, err
	}

	cartEvent := &events.CartUpdatedEvent{
		ID:        cartModel.ID,
		Products:  cartModel.Products,
		Shipping:  cartModel.Shipping,
		Discount:  cartModel.Discount,
		UpdatedAt: cartModel.UpdatedAt,
		Version:   cartModel.Version,
	}

	go cart.cartEventHandler.CartUpdatedEventHandler(ctx, cartEvent)

	return cartModel, nil
}

func (cart *CartCommandHandler) FinalizeCartCommandHandler(ctx context.Context, command *FinalizeCartCommand) (*models.Cart, error) {
	if command.CardNumber == "" {
		return nil, errors.New("card number is required")
	}

	cardNumberEncrypted, kid, err := cart.encryptCardNumber(command.CardNumber)
	if err != nil {
		return nil, errors.New("error encrypt card number")
	}

	cartDto := &dtos.FinalizeCart{
		ID:         command.ID.Hex(),
		CouponID:   command.CouponID.Hex(),
		Products:   command.Products,
		Shipping:   command.Shipping,
		Discount:   command.Discount,
		CardNumber: cardNumberEncrypted,
		Version:    command.Version,
	}

	result := validators.ValidateFinalizeCart(cartDto)
	if result != nil {
		return nil, errors.New(strings.Join(result.([]string), ""))
	}

	err = cart.validateProduct(command.Products)
	if err != nil {
		return nil, err
	}

	cartModel := &models.Cart{
		ID:        command.ID,
		Products:  cartDto.Products,
		Shipping:  cartDto.Shipping,
		Discount:  cartDto.Discount,
		Version:   cartDto.Version,
		UpdatedAt: time.Now().UTC(),
	}

	cartExists, err := cart.cartRepository.FindByID(ctx, cartModel.ID)
	if err != nil {
		return nil, err
	}

	if command.CustomerID != cartExists.CustomerID {
		return nil, errors.New("you can only finalize your cart")
	}

	cartModel, err = cart.cartRepository.Update(ctx, cartModel)
	if err != nil {
		return nil, err
	}

	cartEvent := &events.CartFinalizedEvent{
		ID:         cartModel.ID,
		CustomerID: command.CustomerID,
		Products:   cartModel.Products,
		Sum:        cartModel.SumProduct(),
		Shipping:   cartModel.Shipping,
		Discount:   cartModel.Discount,
		CardNumber: cardNumberEncrypted,
		Kid:        kid,
	}

	go cart.cartEventHandler.CartFinalizedEventHandler(ctx, cartEvent)

	return cartModel, nil
}

func (cart *CartCommandHandler) DeleteCartCommandHandler(ctx context.Context, command *DeleteCartCommand) error {
	cartDto := &dtos.DeleteCart{
		ID: command.ID.Hex(),
	}

	result := validators.ValidateDeleteCart(cartDto)
	if result != nil {
		return errors.New(strings.Join(result.([]string), ""))
	}

	err := cart.cartRepository.Delete(ctx, command.ID)
	if err != nil {
		return err
	}

	cartEvent := &events.CartDeletedEvent{
		ID: command.ID,
	}

	go cart.cartEventHandler.CartDeletedEventHandler(ctx, cartEvent)

	return nil
}

func (cart *CartCommandHandler) validateProduct(products []*models.Product) error {
	for _, product := range products {
		productModel := &dtos.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
		}

		result := validators.ValidateProduct(productModel)
		if result != nil {
			return errors.New(strings.Join(result.([]string), ""))
		}
	}

	return nil
}

func (cart *CartCommandHandler) encryptCardNumber(cartNumber string) ([]byte, string, error) {
	rsaKeys := cart.managerSecurityRSAKeys.GetAllRSAPublicKeys()

	if len(rsaKeys) == 0 {
		return nil, "", nil
	}

	cardNumberEncrypted, err := cart.managerSecurityRSAKeys.Encrypt(cartNumber, rsaKeys[0].Key)
	if err != nil {
		return nil, "", err
	}

	return cardNumberEncrypted, rsaKeys[0].Kid, nil
}
