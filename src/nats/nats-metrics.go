package nats

import (
	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type natsMetric struct {
	config *config.Config
}

var successCartCreated prometheus.Counter
var successCartUpdated prometheus.Counter
var successCartFinalized prometheus.Counter
var successCouponCreated prometheus.Counter
var successCouponUpdated prometheus.Counter
var errorPublish prometheus.Counter

func NewNatsMetric(
	config *config.Config,
) *natsMetric {
	successCartCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_cart_created_total",
			Help: "The total number of success cart created NATS messages",
		},
	)

	successCartUpdated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_cart_updated_total",
			Help: "The total number of success cart updated NATS messages",
		},
	)

	successCartFinalized = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_cart_finalized_total",
			Help: "The total number of success cart finalized NATS messages",
		},
	)

	successCouponCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_coupon_created_total",
			Help: "The total number of success coupon created NATS messages",
		},
	)

	successCouponUpdated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_coupon_updated_total",
			Help: "The total number of success coupon updated NATS messages",
		},
	)

	errorPublish = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_error_publish_message_total",
			Help: "The total number of error NATS publish message",
		},
	)

	return &natsMetric{
		config: config,
	}
}

func (nats *natsMetric) SuccessPublishCartCreated() {
	successCartCreated.Inc()
}

func (nats *natsMetric) SuccessPublishCartUpdated() {
	successCartUpdated.Inc()
}

func (nats *natsMetric) SuccessPublishCartFinalized() {
	successCartFinalized.Inc()
}

func (nats *natsMetric) SuccessPublishCouponCreated() {
	successCouponCreated.Inc()
}

func (nats *natsMetric) SuccessPublishCouponUpdated() {
	successCouponUpdated.Inc()
}

func (nats *natsMetric) ErrorPublish() {
	errorPublish.Inc()
}
