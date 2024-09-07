package v1

import (
	"pyre-promotion/feature-health/v1/controller"

	"github.com/gin-gonic/gin"
)

func HealthV1Route(heatlthRoute *gin.RouterGroup, controller *controller.HealthController) {
	v1Route := heatlthRoute.Group("/v1")

	v1Route.GET("", controller.Health)

}
