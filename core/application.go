package core

import (
	"pyre-promotion/core-internal/conf"
	"pyre-promotion/core-internal/infrastructure"
	health_v1_controller "pyre-promotion/feature-health/v1/controller"
	health_v1_service "pyre-promotion/feature-health/v1/service"

	discount_v1_controller "pyre-promotion/feature-discount/v1/controller"
	discount_v1_service "pyre-promotion/feature-discount/v1/service"

	kafka_producer_controller "pyre-promotion/kafka-produce/controller"
	kafka_producer_service "pyre-promotion/kafka-produce/service"

	"github.com/go-playground/validator/v10"
)

type Application struct {
	HealthV1Service      *health_v1_service.HealthService
	HealthV1Controller *health_v1_controller.HealthController

	KafkaProduceService     *kafka_producer_service.KafkaProduceService
	KafkaProducerController *kafka_producer_controller.KafkaProduceController

	DiscountV1Service    *discount_v1_service.DiscountService
	DiscountV1Controller *discount_v1_controller.DiscountController
}

func NewApplication(env conf.Env) *Application {
	application := &Application{}

	postgresInfra := infrastructure.NewPostgresInfra()
	kafkaInfra := infrastructure.NewKafkaInfra()
	redisInfra := infrastructure.NewRedisInfra()
	middlewareInfra := infrastructure.NewMiddlewareInfra()
	validate := validator.New(validator.WithRequiredStructEnabled())
	productProtoClientInfra := infrastructure.NewProductProtoClientInfra()
	otelInfra := infrastructure.NewOtelInfra()

	application.HealthV1Service = health_v1_service.NewHealthService()
	application.HealthV1Controller = health_v1_controller.NewHealthController(application.HealthV1Service)

	application.DiscountV1Service = discount_v1_service.NewDiscountService(postgresInfra, redisInfra, productProtoClientInfra, otelInfra)
	application.DiscountV1Controller = discount_v1_controller.NewDiscountConttroller(application.DiscountV1Service, validate, middlewareInfra)

	application.KafkaProduceService = kafka_producer_service.NewKafkaProduceService(kafkaInfra.Client)
	application.KafkaProducerController = kafka_producer_controller.NewKafkaProducerController(application.KafkaProduceService, validate, middlewareInfra)

	return application
}
