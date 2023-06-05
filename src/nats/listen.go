package nats

import (
	commands "cart/src/application/commands/cart"
	"cart/src/nats/listeners"

	"github.com/nats-io/nats.go"
	"github.com/oceano-dev/microservices-go-common/config"

	"cart/src/nats/interfaces"

	common_nats "github.com/oceano-dev/microservices-go-common/nats"
	common_service "github.com/oceano-dev/microservices-go-common/services"
)

type listen struct {
	js nats.JetStreamContext
}

const queueGroupName string = "carts-service"

var (
	subscribe          common_nats.Listener
	commandErrorHelper *common_nats.CommandErrorHelper

	cartDeleteCommand *listeners.CartDeleteCommandListener
)

func NewListen(
	config *config.Config,
	js nats.JetStreamContext,
	cartCommandHandler *commands.CartCommandHandler,
	email common_service.EmailService,
	natsMetrics interfaces.NatsMetric,
) *listen {
	subscribe = common_nats.NewListener(js)
	commandErrorHelper = common_nats.NewCommandErrorHelper(config, email)

	cartDeleteCommand = listeners.NewCartDeleteCommandListener(cartCommandHandler, email, commandErrorHelper, natsMetrics)
	return &listen{
		js: js,
	}
}

func (l *listen) Listen() {
	go subscribe.Listener(string(common_nats.OrderCreated), queueGroupName, queueGroupName+"_0", cartDeleteCommand.ProcessCartDeleteCommand())
}
