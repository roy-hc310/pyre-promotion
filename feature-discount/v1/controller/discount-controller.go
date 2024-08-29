package controller

import (
	"net/http"
	"pyre-promotion/core-internal/infrastructure"
	core_model "pyre-promotion/core-internal/model"
	"pyre-promotion/feature-discount/v1/model"
	"pyre-promotion/feature-discount/v1/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DiscountController struct {
	Middleware      *infrastructure.MiddlewareInfra
	Validate        *validator.Validate
	DiscountService *service.DiscountService
}

func NewDiscountConttroller(discountService *service.DiscountService, validate *validator.Validate, middleware *infrastructure.MiddlewareInfra) *DiscountController {
	return &DiscountController{
		Middleware:      middleware,
		Validate:        validate,
		DiscountService: discountService,
	}
}

func (d *DiscountController) CreateDiscount(g *gin.Context) {
	res := core_model.CoreResponseObject{}

	body := model.DiscountRequest{}
	err := g.ShouldBindJSON(&body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = d.Validate.Struct(body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	body.ShopID = d.Middleware.ShopId

	data, statusCode, err := d.DiscountService.CreateDiscount(body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Data = data
	res.Succeeded = true
	g.JSON(statusCode, res)
}

func (d *DiscountController) DetailDiscount(g *gin.Context) {
	res := core_model.CoreResponseObject{}

	id := g.Param("id")

	data, statusCode, err := d.DiscountService.DetailDiscount(id)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Data = data
	res.Succeeded = true
	g.JSON(statusCode, res)
}

func (d *DiscountController) UpdateDiscount(g *gin.Context) {
	res := core_model.CoreResponseObject{}

	id := g.Param("id")
	body := model.DiscountRequest{}
	err := g.ShouldBindJSON(&body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err = d.Validate.Struct(body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	statusCode, err := d.DiscountService.UpdateDiscount(id, body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Succeeded = true
	g.JSON(statusCode, res)
}

func (d *DiscountController) ListDiscounts(g *gin.Context) {
	res := core_model.CoreResponseArray{}

	// query := model.DiscountQuery{}
	query := core_model.CoreQuery{}
	err := g.ShouldBindQuery(&query)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	query.ShopID = d.Middleware.ShopId

	data, statusCode, err := d.DiscountService.ListDiscounts(query)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Data = data
	res.Succeeded = true
	g.JSON(statusCode, res)
}

func (d *DiscountController) DeleteDiscount(g *gin.Context) {
	res := core_model.CoreResponseObject{}

	id := g.Param("id")

	statusCode, err := d.DiscountService.DeleteDiscount(id)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Succeeded = true
	g.JSON(statusCode, res)
}
