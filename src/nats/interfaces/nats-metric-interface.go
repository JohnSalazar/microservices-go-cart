package interfaces

type NatsMetric interface {
	SuccessPublishCartCreated()
	SuccessPublishCartUpdated()
	SuccessPublishCartFinalized()
	SuccessPublishCouponCreated()
	SuccessPublishCouponUpdated()
	ErrorPublish()
}
