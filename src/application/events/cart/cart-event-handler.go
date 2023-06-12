package events

import (
	"context"
	"encoding/json"
	"fmt"

	common_nats "github.com/JohnSalazar/microservices-go-common/nats"
)

type CartEventHandler struct {
	publisher common_nats.Publisher
}

func NewCartEventHandler(
	publisher common_nats.Publisher,
) *CartEventHandler {
	return &CartEventHandler{
		publisher: publisher,
	}
}

func (cart *CartEventHandler) CartCreatedEventHandler(ctx context.Context, event *CartCreatedEvent) error {
	// fmt.Println(event)

	return nil
}

func (cart *CartEventHandler) CartUpdatedEventHandler(ctx context.Context, event *CartUpdatedEvent) error {
	// fmt.Println(event)

	return nil
}

func (cart *CartEventHandler) CartFinalizedEventHandler(ctx context.Context, event *CartFinalizedEvent) error {
	data, _ := json.Marshal(event)
	err := cart.publisher.Publish(string(common_nats.OrderCreate), data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (cart *CartEventHandler) CartDeletedEventHandler(ctx context.Context, event *CartDeletedEvent) error {
	// fmt.Println(event)

	return nil
}
