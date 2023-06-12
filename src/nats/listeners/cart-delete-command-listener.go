package listeners

import (
	commands "cart/src/application/commands/cart"
	"context"
	"encoding/json"
	"fmt"
	"log"

	natsMetrics "cart/src/nats/interfaces"

	common_nats "github.com/JohnSalazar/microservices-go-common/nats"
	common_service "github.com/JohnSalazar/microservices-go-common/services"
	trace "github.com/JohnSalazar/microservices-go-common/trace/otel"
	"github.com/nats-io/nats.go"
)

type CartDeleteCommandListener struct {
	commandHandler *commands.CartCommandHandler
	email          common_service.EmailService
	errorHelper    *common_nats.CommandErrorHelper
	natsMetrics    natsMetrics.NatsMetric
}

func NewCartDeleteCommandListener(
	commandHandler *commands.CartCommandHandler,
	email common_service.EmailService,
	errorHelper *common_nats.CommandErrorHelper,
	natsMetrics natsMetrics.NatsMetric,
) *CartDeleteCommandListener {
	return &CartDeleteCommandListener{
		commandHandler: commandHandler,
		email:          email,
		errorHelper:    errorHelper,
		natsMetrics:    natsMetrics,
	}
}

func (c *CartDeleteCommandListener) ProcessCartDeleteCommand() nats.MsgHandler {
	return func(msg *nats.Msg) {
		ctx := context.Background()
		_, span := trace.NewSpan(ctx, fmt.Sprintf("publish.%s\n", msg.Subject))
		defer span.End()

		cartCommand := &commands.DeleteCartCommand{}
		err := json.Unmarshal(msg.Data, cartCommand)
		if c.errorHelper.CheckUnmarshal(msg, err) == nil {
			c.natsMetrics.ErrorPublish()
			err = c.commandHandler.DeleteCartCommandHandler(ctx, cartCommand)
			c.errorHelper.CheckCommandError(span, msg, err)
		}

		err = msg.Ack()
		if err != nil {
			log.Printf("stan msg.Ack error: %v\n", err)
		}
	}
}
