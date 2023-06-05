package controllers

import (
	commands "cart/src/application/commands/coupon"
	"cart/src/repositories/interfaces"
	"net/http"
	"strconv"

	natsMetrics "cart/src/nats/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/oceano-dev/microservices-go-common/helpers"
	"github.com/oceano-dev/microservices-go-common/httputil"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
)

type CouponController struct {
	couponRepository     interfaces.CouponRepository
	couponCommandHandler *commands.CouponCommandHandler
	natsMetrics          natsMetrics.NatsMetric
}

func NewCouponController(
	couponRepository interfaces.CouponRepository,
	couponCommandHandler *commands.CouponCommandHandler,
	natsMetrics natsMetrics.NatsMetric,
) *CouponController {
	return &CouponController{
		couponRepository:     couponRepository,
		couponCommandHandler: couponCommandHandler,
		natsMetrics:          natsMetrics,
	}
}

func (coupon *CouponController) GetAll(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CouponController.GetAll")
	defer span.End()

	name := c.Param("name")

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "page is required")
		return
	}

	size, err := strconv.Atoi(c.Param("size"))
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "size is required")
		return
	}

	_coupon, err := coupon.couponRepository.GetAll(c.Request.Context(), name, page, size)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "coupons not found")
		return
	}

	c.JSON(http.StatusOK, _coupon)
}

func (coupon *CouponController) Get(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CouponController.Get")
	defer span.End()

	isID := helpers.IsValidID(c.Param("id"))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	ID := helpers.StringToID(c.Param("id"))

	_coupon, err := coupon.couponRepository.FindByID(c.Request.Context(), ID)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "coupon not found")
		return
	}

	c.JSON(http.StatusOK, _coupon)
}

func (coupon *CouponController) GetByName(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CouponController.GetByName")
	defer span.End()

	name := c.Param("name")
	if len(name) == 0 {
		httputil.NewResponseError(c, http.StatusBadRequest, "coupon name invalid")
		return
	}

	_coupon, err := coupon.couponRepository.GetByName(c.Request.Context(), name)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "coupon unavailable")
		return
	}

	c.JSON(http.StatusOK, _coupon)
}

func (coupon *CouponController) Create(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "CouponController.Create")
	defer span.End()

	createCouponCommand := &commands.CreateCouponCommand{}
	err := c.BindJSON(createCouponCommand)
	if err != nil {
		trace.FailSpan(span, "Error json parse")
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	couponModel, err := coupon.couponCommandHandler.CreateCouponCommandHandler(ctx, createCouponCommand)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	coupon.natsMetrics.SuccessPublishCouponCreated()

	c.JSON(http.StatusCreated, couponModel)
}

func (coupon *CouponController) Update(c *gin.Context) {
	ctx, span := trace.NewSpan(c.Request.Context(), "CouponController.Update")
	defer span.End()

	updateCouponCommand := &commands.UpdateCouponCommand{}
	err := c.BindJSON(updateCouponCommand)
	if err != nil {
		trace.FailSpan(span, "Error json parse")
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	couponModel, err := coupon.couponCommandHandler.UpdateCouponCommandHandler(ctx, updateCouponCommand)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	coupon.natsMetrics.SuccessPublishCouponUpdated()

	c.JSON(http.StatusOK, couponModel)
}
