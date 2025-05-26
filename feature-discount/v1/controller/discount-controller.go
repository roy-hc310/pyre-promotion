package controller

import (
	"context"
	"fmt"
	"net/http"
	core_model "pyre-promotion/core-internal/model"
	"pyre-promotion/core-internal/utils"
	"pyre-promotion/feature-discount/v1/model"
	"pyre-promotion/feature-discount/v1/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type DiscountController struct {
	Validate        *validator.Validate
	DiscountService *service.DiscountService
}

func NewDiscountConttroller(discountService *service.DiscountService, validate *validator.Validate) *DiscountController {
	return &DiscountController{
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

	shopId, _ := g.Get(utils.XShopId)
	body.ShopID = fmt.Sprintf("%v", shopId)

	data, traceID, statusCode, err := d.DiscountService.CreateDiscount(g.Request.Context(), body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Data = data
	res.TraceID = traceID
	res.Succeeded = true
	g.JSON(statusCode, res)
}

func (d *DiscountController) DetailDiscount(g *gin.Context) {
	res := core_model.CoreResponseObject{}

	id := g.Param("id")

	data, traceID, statusCode, err := d.DiscountService.DetailDiscount(g.Request.Context(), id)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Data = data
	res.TraceID = traceID
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

	traceID, statusCode, err := d.DiscountService.UpdateDiscount(g.Request.Context(), id, body)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.TraceID = traceID
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

	shopId, _ := g.Get(utils.XShopId)
	query.ShopID = fmt.Sprintf("%v", shopId)

	data, pagination, traceID, statusCode, err := d.DiscountService.ListDiscounts(g.Request.Context(), query)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.Data = data
	res.TraceID = traceID
	res.Succeeded = true
	res.Pagination = pagination
	g.JSON(statusCode, res)
}

func (d *DiscountController) DeleteDiscount(g *gin.Context) {
	ctx, _ := context.WithTimeout(g.Request.Context(), utils.DefaultContextTimeOut)
	res := core_model.CoreResponseObject{}

	id := g.Param("id")

	traceID, statusCode, err := d.DiscountService.DeleteDiscount(ctx, id)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		g.AbortWithStatusJSON(statusCode, res)
		return
	}

	res.TraceID = traceID
	res.Succeeded = true
	g.JSON(statusCode, res)
}
