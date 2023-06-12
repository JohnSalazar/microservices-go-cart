package events

import (
	"context"

	common_nats "github.com/JohnSalazar/microservices-go-common/nats"
)

type CouponEventHandler struct {
	publisher common_nats.Publisher
}

func NewCouponEventHandler(
	publisher common_nats.Publisher,
) *CouponEventHandler {
	return &CouponEventHandler{
		publisher: publisher,
	}
}

func (coupon *CouponEventHandler) CouponCreatedEventHandler(ctx context.Context, event *CouponCreatedEvent) error {
	// fmt.Println(event)

	return nil
}

func (coupon *CouponEventHandler) CouponUpdatedEventHandler(ctx context.Context, event *CouponUpdatedEvent) error {
	// fmt.Println(event)

	return nil
}
