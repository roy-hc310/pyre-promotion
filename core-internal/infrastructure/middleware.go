package infrastructure

import (
	"net/http"
	"pyre-promotion/core-internal/utils"

	"github.com/gin-gonic/gin"
)

type MiddlewareInfra struct {
	ShopId   string
	Language string
}

func NewMiddlewareInfra() *MiddlewareInfra {
	return &MiddlewareInfra{}
}

func (m *MiddlewareInfra) Prepare(g *gin.Context) {
	shopId := g.Request.Header.Get(utils.XShopId)
	if shopId == "" {
		g.AbortWithStatusJSON(http.StatusUnauthorized, "Missing x-shop-id")
		return
	}

	m.ShopId = shopId
	g.Next()
}
