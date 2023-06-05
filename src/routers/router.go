package routers

import (
	"cart/src/controllers"
	"fmt"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	common_service "github.com/oceano-dev/microservices-go-common/services"
)

type router struct {
	config           *config.Config
	serviceMetrics   common_service.Metrics
	authentication   *middlewares.Authentication
	cartController   *controllers.CartController
	couponController *controllers.CouponController
}

func NewRouter(
	config *config.Config,
	serviceMetrics common_service.Metrics,
	authentication *middlewares.Authentication,
	cartController *controllers.CartController,
	couponController *controllers.CouponController,
) *router {
	return &router{
		config:           config,
		serviceMetrics:   serviceMetrics,
		authentication:   authentication,
		cartController:   cartController,
		couponController: couponController,
	}
}

func (r *router) RouterSetup() *gin.Engine {
	router := r.initRouter()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())
	router.Use(location.Default())
	router.Use(otelgin.Middleware(r.config.Jaeger.ServiceName))
	router.Use(middlewares.Metrics(r.serviceMetrics))

	router.GET("/healthy", middlewares.Healthy())
	router.GET("/metrics", middlewares.MetricsHandler())

	v1 := router.Group(fmt.Sprintf("/api/%s", r.config.ApiVersion))

	v1.GET("/", r.authentication.Verify(),
		r.cartController.Get)
	v1.POST("/", r.authentication.Verify(),
		r.cartController.Create)
	v1.PUT("/:id", r.authentication.Verify(),
		r.cartController.Update)
	v1.PUT("/finalize/:id", r.authentication.Verify(),
		r.cartController.Finalize)

	v1.GET("/coupons/:name/:page/:size", r.authentication.Verify(),
		r.couponController.GetAll)
	v1.GET("/coupon/:id", r.authentication.Verify(),
		r.couponController.Get)
	v1.GET("/coupon/name/:name", r.authentication.Verify(),
		r.couponController.GetByName)
	v1.POST("/coupon", r.authentication.Verify(),
		middlewares.Authorization("admin", "create,read"),
		r.couponController.Create)
	v1.PUT("/coupon/:id", r.authentication.Verify(),
		middlewares.Authorization("admin", "create,read"),
		r.couponController.Update)

	return router
}

func (r *router) initRouter() *gin.Engine {
	if r.config.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	return gin.New()
}
