package v1

import (
	"pyre-promotion/core-internal/utils"
	"pyre-promotion/feature-discount/v1/controller"

	"github.com/gin-gonic/gin"
)

func DiscountV1Route(discountRoute *gin.RouterGroup, controller *controller.DiscountController) {
	v1Route := discountRoute.Group("/v1")

	v1Route.POST("", utils.Middleware, controller.CreateDiscount)
	v1Route.GET("/:id", utils.Middleware, controller.DetailDiscount)
	v1Route.PUT("/:id", utils.Middleware, controller.UpdateDiscount)
	v1Route.GET("", utils.Middleware, controller.ListDiscounts)
	v1Route.DELETE("/:id", utils.Middleware, controller.DeleteDiscount)
}
