package core

import (
	discount_v1 "pyre-promotion/feature-discount/v1"
	health_v1 "pyre-promotion/feature-health/v1"
	kafka_producer "pyre-promotion/kafka-produce"

	"github.com/gin-gonic/gin"
)

func Router(g *gin.Engine, application *Application) error {

	route := g.Group("/api")

	healthRoute := route.Group("/health")
	health_v1.HealthV1Route(healthRoute, application.HealthV1Controller)

	KafkaProduceRoute := route.Group("/kafka")
	kafka_producer.KafkaProduceRoute(KafkaProduceRoute, application.KafkaProducerController)

	discountRoute := route.Group("/discount")
	discount_v1.DiscountV1Route(discountRoute, application.DiscountV1Controller)

	return nil
}
