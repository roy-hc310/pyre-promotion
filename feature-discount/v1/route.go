package v1

import (
	"pyre-promotion/feature-discount/v1/controller"

	"github.com/gin-gonic/gin"
)

func DiscountV1Route(discountRoute *gin.RouterGroup, controller *controller.DiscountController) {
	v1Route := discountRoute.Group("/v1")

	v1Route.POST("", controller.Middleware.Prepare, controller.CreateDiscount)
	v1Route.GET("/:id", controller.Middleware.Prepare, controller.DetailDiscount)
	v1Route.PUT("/:id", controller.Middleware.Prepare, controller.UpdateDiscount)
	v1Route.GET("", controller.Middleware.Prepare, controller.ListDiscounts)
	v1Route.DELETE("/:id", controller.Middleware.Prepare, controller.DeleteDiscount)
}
