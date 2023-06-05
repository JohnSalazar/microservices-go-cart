package controllers

import (
	"context"
	"errors"
	"net/http"

	commands_cart "cart/src/application/commands/cart"
	commands_coupon "cart/src/application/commands/coupon"
	"cart/src/repositories/interfaces"

	natsMetrics "cart/src/nats/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/oceano-dev/microservices-go-common/helpers"
	"github.com/oceano-dev/microservices-go-common/httputil"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartController struct {
	cartRepository       interfaces.CartRepository
	couponRepository     interfaces.CouponRepository
	cartCommandHandler   *commands_cart.CartCommandHandler
	couponCommandHandler *commands_coupon.CouponCommandHandler
	natsMetrics          natsMetrics.NatsMetric
}

func NewCartController(
	cartRepository interfaces.CartRepository,
	couponRepository interfaces.CouponRepository,
	cartCommandHandler *commands_cart.CartCommandHandler,
	couponCommandHandler *commands_coupon.CouponCommandHandler,
	natsMetrics natsMetrics.NatsMetric,
) *CartController {
	return &CartController{
		cartRepository:       cartRepository,
		couponRepository:     couponRepository,
		cartCommandHandler:   cartCommandHandler,
		couponCommandHandler: couponCommandHandler,
		natsMetrics:          natsMetrics,
	}
}

func (cart *CartController) Get(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CartController.Get")
	defer span.End()

	ID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid customer")
		return
	}

	isID := helpers.IsValidID(ID.(string))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid customerId")
		return
	}

	customerID := helpers.StringToID(ID.(string))

	cartModel, err := cart.cartRepository.FindByCustomerID(c.Request.Context(), customerID)
	if cartModel == nil || err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "cart not found")
		return
	}

	// coupon := &models.Coupon{
	// 	Name:         "Discount 10%",
	// 	Value:        10,
	// 	IsPercentage: true,
	// }
	// fmt.Printf("%s: %v\n", coupon.Name, cartModel.ApplyCoupon(coupon))

	// fmt.Printf("Total: %v\n", cartModel.Total())

	// coupon.Name = "Discount $100"
	// coupon.Value = 100
	// coupon.IsPercentage = false
	// fmt.Printf("%s: %v\n", coupon.Name, cartModel.ApplyCoupon(coupon))

	c.JSON(http.StatusOK, cartModel)
}

func (cart *CartController) Create(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "CartController.Create")
	defer span.End()

	customerID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid customer")
		return
	}

	ID := helpers.StringToID(customerID.(string))

	createCartCommand := &commands_cart.CreateCartCommand{}
	err := c.BindJSON(createCartCommand)
	if err != nil {
		trace.FailSpan(span, "Error json parse")
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	createCartCommand.CustomerID = ID

	cartModel, err := cart.cartCommandHandler.CreateCartCommandHandler(ctx, createCartCommand)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	cart.natsMetrics.SuccessPublishCartCreated()

	c.JSON(http.StatusCreated, cartModel)
}

func (cart *CartController) Update(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "CartController.Update")
	defer span.End()

	isID := helpers.IsValidID(c.Param("id"))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid Id")
		return
	}

	ID := helpers.StringToID(c.Param("id"))

	customerID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid customer")
		return
	}

	updateCartCommand := &commands_cart.UpdateCartCommand{}
	err := c.BindJSON(updateCartCommand)
	if err != nil {
		trace.FailSpan(span, "Error json parse")
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if ID != updateCartCommand.ID {
		httputil.NewResponseError(c, http.StatusForbidden, "you can only change your cart")
		return
	}

	updateCartCommand.CustomerID = helpers.StringToID(customerID.(string))

	cartModel, err := cart.cartCommandHandler.UpdateCartCommandHandler(ctx, updateCartCommand)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "cart not found")
		return
	}

	cart.natsMetrics.SuccessPublishCartUpdated()

	c.JSON(http.StatusOK, cartModel)
}

func (cart *CartController) Finalize(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "CartController.Finalize")
	defer span.End()

	isID := helpers.IsValidID(c.Param("id"))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid Id")
		return
	}

	ID := helpers.StringToID(c.Param("id"))

	customerID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid customer")
		return
	}

	finalizeCartCommand := &commands_cart.FinalizeCartCommand{}
	err := c.BindJSON(finalizeCartCommand)
	if err != nil {
		trace.FailSpan(span, "Error json parse")
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	if ID != finalizeCartCommand.ID {
		httputil.NewResponseError(c, http.StatusForbidden, "you can only finalize your cart")
		return
	}

	_cart, err := cart.cartRepository.FindByID(ctx, ID)
	if err != nil || _cart == nil || _cart.Version != finalizeCartCommand.Version {
		httputil.NewResponseError(c, http.StatusBadRequest, "cart not found")
		return
	}

	finalizeCartCommand.CustomerID = helpers.StringToID(customerID.(string))

	if !finalizeCartCommand.CouponID.IsZero() {
		err := cart.updateCoupon(ctx, finalizeCartCommand.CouponID)
		if err != nil {
			httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
			return
		}
	}

	cartModel, err := cart.cartCommandHandler.FinalizeCartCommandHandler(ctx, finalizeCartCommand)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	cart.natsMetrics.SuccessPublishCartFinalized()

	c.JSON(http.StatusOK, cartModel)
}

func (cart *CartController) updateCoupon(ctx context.Context, couponID primitive.ObjectID) error {
	coupon, err := cart.couponRepository.FindByID(ctx, couponID)
	if coupon == nil || err != nil {
		return errors.New("coupon not found")
	}

	if !coupon.Active || coupon.Quantity == 0 {
		return errors.New("coupon unavailable")
	}

	coupon.Quantity--
	updateCouponCommand := &commands_coupon.UpdateCouponCommand{
		ID:           coupon.ID,
		Name:         coupon.Name,
		Value:        coupon.Value,
		IsPercentage: coupon.IsPercentage,
		Quantity:     coupon.Quantity,
		Active:       coupon.Active,
		Version:      coupon.Version,
	}

	_, err = cart.couponCommandHandler.UpdateCouponCommandHandler(ctx, updateCouponCommand)
	if err != nil {
		return err
	}

	return nil
}
