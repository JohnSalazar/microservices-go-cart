package commands

import (
	events "cart/src/application/events/coupon"
	"cart/src/dtos"
	"cart/src/models"
	"cart/src/repositories/interfaces"
	"cart/src/validators"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CouponCommandHandler struct {
	couponRepository   interfaces.CouponRepository
	couponEventHandler *events.CouponEventHandler
}

func NewCouponCommandHandler(
	couponRepository interfaces.CouponRepository,
	couponEventHandler *events.CouponEventHandler,
) *CouponCommandHandler {
	return &CouponCommandHandler{
		couponRepository:   couponRepository,
		couponEventHandler: couponEventHandler,
	}
}

func (coupon *CouponCommandHandler) CreateCouponCommandHandler(ctx context.Context, command *CreateCouponCommand) (*models.Coupon, error) {
	couponDto := &dtos.AddCoupon{
		Name:         command.Name,
		Value:        command.Value,
		IsPercentage: command.IsPercentage,
		Quantity:     command.Quantity,
		Active:       command.Active,
	}

	result := validators.ValidateAddCoupon(couponDto)
	if result != nil {
		return nil, errors.New(strings.Join(result.([]string), ""))
	}

	couponModel := &models.Coupon{
		ID:           primitive.NewObjectID(),
		Name:         couponDto.Name,
		Value:        couponDto.Value,
		IsPercentage: couponDto.IsPercentage,
		Quantity:     couponDto.Quantity,
		Active:       couponDto.Active,
		CreatedAt:    time.Now().UTC(),
	}

	couponExists, _ := coupon.couponRepository.FindByName(ctx, couponDto.Name)
	if couponExists != nil {
		return nil, errors.New("already a coupon with this name")
	}

	couponModel, err := coupon.couponRepository.Create(ctx, couponModel)
	if err != nil {
		return nil, err
	}

	couponEvent := &events.CouponCreatedEvent{
		ID:           couponModel.ID,
		Name:         couponModel.Name,
		Value:        couponModel.Value,
		IsPercentage: couponModel.IsPercentage,
		Quantity:     couponModel.Quantity,
		Active:       couponModel.Active,
		CreatedAt:    couponModel.CreatedAt,
		UpdatedAt:    couponModel.UpdatedAt,
		Version:      couponModel.Version,
		Deleted:      couponModel.Deleted,
	}

	go coupon.couponEventHandler.CouponCreatedEventHandler(ctx, couponEvent)

	return couponModel, nil
}

func (coupon *CouponCommandHandler) UpdateCouponCommandHandler(ctx context.Context, command *UpdateCouponCommand) (*models.Coupon, error) {
	couponDto := *&dtos.UpdateCoupon{
		ID:           command.ID.Hex(),
		Name:         command.Name,
		Value:        command.Value,
		IsPercentage: command.IsPercentage,
		Quantity:     command.Quantity,
		Active:       command.Active,
		Version:      command.Version,
	}

	result := validators.ValidateUpdateCoupon(&couponDto)
	if result != nil {
		return nil, errors.New(strings.Join(result.([]string), ""))
	}

	couponModel := &models.Coupon{
		ID:           command.ID,
		Name:         couponDto.Name,
		Value:        couponDto.Value,
		IsPercentage: couponDto.IsPercentage,
		Quantity:     couponDto.Quantity,
		Active:       couponDto.Active,
		Version:      couponDto.Version,
		UpdatedAt:    time.Now().UTC(),
	}

	couponExists, _ := coupon.couponRepository.FindByName(ctx, couponDto.Name)
	if couponExists != nil && command.ID != couponExists.ID {
		return nil, errors.New("already a coupon with this name")
	}

	couponModel, err := coupon.couponRepository.Update(ctx, couponModel)
	if err != nil {
		return nil, err
	}

	couponEvent := &events.CouponUpdatedEvent{
		ID:           couponModel.ID,
		Name:         couponModel.Name,
		Value:        couponModel.Value,
		IsPercentage: couponModel.IsPercentage,
		Quantity:     couponModel.Quantity,
		Active:       couponModel.Active,
		UpdatedAt:    couponModel.UpdatedAt,
	}

	go coupon.couponEventHandler.CouponUpdatedEventHandler(ctx, couponEvent)

	return couponModel, nil
}
