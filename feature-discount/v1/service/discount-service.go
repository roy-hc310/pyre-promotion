package service

import (
	"encoding/json"
	"errors"
	"pyre-promotion/core-internal/infrastructure"
	core_model "pyre-promotion/core-internal/model"
	"pyre-promotion/core-internal/utils"
	"pyre-promotion/feature-discount/v1/model"
	"strings"
	"sync"
	"time"

	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DiscountService struct {
	PostgresInfra           *infrastructure.PostgresInfra
	RedisInfra              *infrastructure.RedisInfra
	ProductProtoClientInfra *infrastructure.ProductProtoClientInfra
	OtelInfra               *infrastructure.OtelInfra
}

func NewDiscountService(postgresInfra *infrastructure.PostgresInfra, redisInfra *infrastructure.RedisInfra, productProtoClientInfra *infrastructure.ProductProtoClientInfra, otelInfra *infrastructure.OtelInfra) *DiscountService {
	return &DiscountService{
		PostgresInfra:           postgresInfra,
		RedisInfra:              redisInfra,
		ProductProtoClientInfra: productProtoClientInfra,
		OtelInfra:               otelInfra,
	}
}

func (d *DiscountService) CreateDiscount(ctx context.Context, data model.DiscountRequest) (res core_model.CoreIdResponse, traceID string, statusCode int, err error) {
	ctx, span := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount")
	traceID = span.SpanContext().TraceID().String()

	_, validationSpan := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount-Validation")
	productIds := []string{}
	for i := 0; i < len(data.Products); i++ {
		productIds = append(productIds, data.Products[i].SKU)
	}

	if utils.GlobalEnv.Debugging == false {
		success, err := d.ProductProtoClientInfra.GetProduct(productIds)
		if err != nil {
			return res, traceID, http.StatusInternalServerError, err
		}

		if !success.Success {
			return res, traceID, http.StatusBadRequest, errors.New("product not found")
		}
	}

	validationSpan.End()

	// promotions
	promotionUUID := uuid.New().String()
	data.UUID = promotionUUID
	data.PromotionType = utils.PromotionTypeDiscount
	promotionsInterface := [][]interface{}{
		{data.UUID, data.Name, data.PromotionType, data.Code, data.StartTime, data.EndTime, data.ShopID, data.UsageQuantity, data.UsageLimitPerUser},
	}

	promotionString := utils.PrepareInsertQuery(utils.PromotionTableName, utils.PromotionColumnsListForInsert, promotionsInterface)

	var promotionParams []interface{}
	for _, promotion := range promotionsInterface {
		promotionParams = append(promotionParams, promotion...)
	}

	_, createTxSpan := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount-Tx")

	tx, err := d.PostgresInfra.DbWritePool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}
	defer tx.Rollback(ctx)

	createTxSpan.End()

	_, dbSpan := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount-DB")
	dbSpan.AddEvent("insert promotions", trace.WithAttributes(
		attribute.String("query", promotionString),
		attribute.String("params", fmt.Sprintf("%v", promotionParams)),
	))

	_, err = tx.Exec(ctx, promotionString, promotionParams...)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	// products
	var productsInterface [][]interface{}
	var variantsInterface [][]interface{}

	for i := 0; i < len(data.Products); i++ {
		data.Products[i].UUID = uuid.New().String()
		data.Products[i].PromotionID = promotionUUID
		productInterface := []interface{}{
			data.Products[i].UUID,
			data.Products[i].PromotionID,
			data.Products[i].SKU,
			data.Products[i].Name,
			data.Products[i].PurchaseLimit,
		}
		productsInterface = append(productsInterface, productInterface)

		for j := 0; j < len(data.Products[i].ProductVariants); j++ {
			data.Products[i].ProductVariants[j].UUID = uuid.New().String()
			data.Products[i].ProductVariants[j].PromotionID = promotionUUID
			data.Products[i].ProductVariants[j].ProductID = data.Products[i].UUID
			variantInterface := []interface{}{
				data.Products[i].ProductVariants[j].UUID,
				data.Products[i].ProductVariants[j].PromotionID,
				data.Products[i].ProductVariants[j].ProductID,
				data.Products[i].ProductVariants[j].SKU,
				data.Products[i].ProductVariants[j].Name,
				data.Products[i].ProductVariants[j].DiscountedPrice,
				data.Products[i].ProductVariants[j].DiscountedPercentage,
				data.Products[i].ProductVariants[j].StockLimit,
				data.Products[i].ProductVariants[j].IsActive,
			}
			variantsInterface = append(variantsInterface, variantInterface)
		}
	}

	productString := utils.PrepareInsertQuery(utils.ProductTableName, utils.ProductColumnsListForInsert, productsInterface)

	var productParams []interface{}
	for _, product := range productsInterface {
		productParams = append(productParams, product...)
	}

	dbSpan.AddEvent("insert products", trace.WithAttributes(
		attribute.String("query", productString),
		attribute.String("params", fmt.Sprintf("%v", productParams)),
	))

	_, err = tx.Exec(ctx, productString, productParams...)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	// variants
	variantString := utils.PrepareInsertQuery(utils.VariantTableName, utils.VariantColumnsListForInsert, variantsInterface)

	var variantParams []interface{}
	for _, variant := range variantsInterface {
		variantParams = append(variantParams, variant...)
	}

	dbSpan.AddEvent("insert variants", trace.WithAttributes(
		attribute.String("query", variantString),
		attribute.String("params", fmt.Sprintf("%v", variantParams)),
	))

	_, err = tx.Exec(ctx, variantString, variantParams...)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	dbSpan.End()
	span.End()
	res.Id = promotionUUID

	return res, traceID, http.StatusOK, nil
}

func (d *DiscountService) CreateBulkDiscount(ctx context.Context, data []model.DiscountRequest) (res []core_model.CoreIdResponse, traceID string, statusCode int, err error) {
	ctx, span := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount")
	traceID = span.SpanContext().TraceID().String()

	response := []core_model.CoreIdResponse{}

	// promotions
	var promotionsInterface [][]interface{}
	var productsInterface [][]interface{}
	var variantsInterface [][]interface{}

	for i := 0; i < len(data); i++ {
		promotionUUID := uuid.New().String()
		data[i].UUID = promotionUUID
		data[i].PromotionType = utils.PromotionTypeDiscount
		response = append(response, core_model.CoreIdResponse{Id: promotionUUID})
		promotionInterface := []interface{}{
			data[i].UUID, data[i].Name, data[i].PromotionType, data[i].Code, data[i].StartTime, data[i].EndTime, data[i].ShopID, data[i].UsageQuantity, data[i].UsageLimitPerUser,
		}
		promotionsInterface = append(promotionsInterface, promotionInterface)

		for j := 0; j < len(data[i].Products); j++ {
			data[i].Products[j].UUID = uuid.New().String()
			data[i].Products[j].PromotionID = promotionUUID
			productInterface := []interface{}{
				data[i].Products[j].UUID,
				data[i].Products[j].PromotionID,
				data[i].Products[j].SKU,
				data[i].Products[j].Name,
				data[i].Products[j].PurchaseLimit,
			}
			productsInterface = append(productsInterface, productInterface)

			for k := 0; k < len(data[i].Products[j].ProductVariants); k++ {
				data[i].Products[j].ProductVariants[k].UUID = uuid.New().String()
				data[i].Products[j].ProductVariants[k].PromotionID = promotionUUID
				data[i].Products[j].ProductVariants[k].ProductID = data[i].Products[j].UUID
				variantInterface := []interface{}{
					data[i].Products[j].ProductVariants[k].UUID,
					data[i].Products[j].ProductVariants[k].PromotionID,
					data[i].Products[j].ProductVariants[k].ProductID,
					data[i].Products[j].ProductVariants[k].SKU,
					data[i].Products[j].ProductVariants[k].Name,
					data[i].Products[j].ProductVariants[k].DiscountedPrice,
					data[i].Products[j].ProductVariants[k].DiscountedPercentage,
					data[i].Products[j].ProductVariants[k].StockLimit,
					data[i].Products[j].ProductVariants[k].IsActive,
				}
				variantsInterface = append(variantsInterface, variantInterface)
			}
		}
	}

	//promotions
	promotionString := utils.PrepareInsertQuery(utils.PromotionTableName, utils.PromotionColumnsListForInsert, promotionsInterface)
	var promotionParams []interface{}
	for _, promotion := range promotionsInterface {
		promotionParams = append(promotionParams, promotion...)
	}

	_, createTxSpan := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount-Tx")

	tx, err := d.PostgresInfra.DbWritePool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}
	defer tx.Rollback(ctx)

	createTxSpan.End()

	_, dbSpan := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount-DB")
	dbSpan.AddEvent("insert promotions", trace.WithAttributes(
		attribute.String("query", promotionString),
		attribute.String("params", fmt.Sprintf("%v", promotionParams)),
	))

	_, err = tx.Exec(ctx, promotionString, promotionParams...)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	// products
	productString := utils.PrepareInsertQuery(utils.ProductTableName, utils.ProductColumnsListForInsert, productsInterface)
	var productParams []interface{}
	for _, product := range productsInterface {
		productParams = append(productParams, product...)
	}

	dbSpan.AddEvent("insert products", trace.WithAttributes(
		attribute.String("query", productString),
		attribute.String("params", fmt.Sprintf("%v", productParams)),
	))

	_, err = tx.Exec(ctx, productString, productParams...)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	// variants
	variantString := utils.PrepareInsertQuery(utils.VariantTableName, utils.VariantColumnsListForInsert, variantsInterface)
	var variantParams []interface{}
	for _, variant := range variantsInterface {
		variantParams = append(variantParams, variant...)
	}

	dbSpan.AddEvent("insert variants", trace.WithAttributes(
		attribute.String("query", variantString),
		attribute.String("params", fmt.Sprintf("%v", variantParams)),
	))

	_, err = tx.Exec(ctx, variantString, variantParams...)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}

	dbSpan.End()
	span.End()

	res = response

	return res, traceID, http.StatusOK, nil
}

func (d *DiscountService) DetailDiscount(ctx context.Context, id string) (res model.DiscountResponse, traceID string, statusCode int, err error) {
	ctx, span := d.OtelInfra.Tracer.Start(ctx, "ListDiscounts")
	traceID = span.SpanContext().TraceID().String()

	promotionString := `SELECT p.id, p.created_at, p.updated_at, p.deleted_at, p.uuid, p.name, p.promotion_type, p.code, p.start_time, p.end_time, p.shop_id, p.usage_quantity, p.usage_limit_per_user, pr.id, pr.uuid, pv.promotion_id, pr.sku, pr.name, pr.purchase_limit,pv.id, pv.uuid, pv.promotion_id, pv.product_id, pv.sku, pv.name, pv.discounted_price, pv.discounted_percentage, pv.stock_limit, pv.is_active FROM promotion.promotions p LEFT JOIN promotion.products pr on pr.promotion_id = p.uuid LEFT JOIN promotion.product_variants pv on pv.product_id = pr.uuid WHERE p.uuid = $1 AND p.deleted_at IS NULL;`

	cache, err := d.RedisInfra.Client.Get(promotionString).Bytes()
	if err == nil {
		err = json.Unmarshal(cache, &res)
		if err == nil {
			return res, traceID, http.StatusOK, nil
		}
	}

	_, dbSpan := d.OtelInfra.Tracer.Start(ctx, "DetailDiscount-DB")
	dbSpan.AddEvent("get promotion", trace.WithAttributes(
		attribute.String("query", promotionString),
		attribute.String("params", fmt.Sprintf("%v", id)),
	))

	rows, err := d.PostgresInfra.DbReadPool.Query(ctx, promotionString, id)
	if err != nil {
		return res, traceID, http.StatusInternalServerError, err
	}
	defer rows.Close()

	dbSpan.End()

	products := []core_model.CoreProduct{}
	promotion := core_model.CorePromotion{}
	mapProduct := make(map[string]core_model.CoreProduct)
	mapProductIdVariants := make(map[string][]core_model.CoreProductVariant)

	for rows.Next() {
		var rawPromotion core_model.CorePromotion
		var rawProduct core_model.CoreProduct
		var rawProductVariant core_model.CoreProductVariant

		err = rows.Scan(
			&rawPromotion.ID,
			&rawPromotion.CreatedAt,
			&rawPromotion.UpdatedAt,
			&rawPromotion.DeletedAt,
			&rawPromotion.UUID,
			&rawPromotion.Name,
			&rawPromotion.PromotionType,
			&rawPromotion.Code,
			&rawPromotion.StartTime,
			&rawPromotion.EndTime,
			&rawPromotion.ShopID,
			&rawPromotion.UsageQuantity,
			&rawPromotion.UsageLimitPerUser,

			&rawProduct.ID,
			&rawProduct.UUID,
			&rawProduct.PromotionID,
			&rawProduct.SKU,
			&rawProduct.Name,
			&rawProduct.PurchaseLimit,

			&rawProductVariant.ID,
			&rawProductVariant.UUID,
			&rawProductVariant.PromotionID,
			&rawProductVariant.ProductID,
			&rawProductVariant.SKU,
			&rawProductVariant.Name,
			&rawProductVariant.DiscountedPrice,
			&rawProductVariant.DiscountedPercentage,
			&rawProductVariant.StockLimit,
			&rawProductVariant.IsActive)
		if err != nil {
			return res, traceID, http.StatusInternalServerError, err
		}

		if _, ok := mapProductIdVariants[rawProductVariant.ProductID]; !ok {
			mapProductIdVariants[rawProductVariant.ProductID] = []core_model.CoreProductVariant{}
		}
		mapProductIdVariants[rawProductVariant.ProductID] = append(mapProductIdVariants[rawProductVariant.ProductID], rawProductVariant)

		if _, ok := mapProduct[rawProduct.UUID]; !ok {
			mapProduct[rawProduct.UUID] = rawProduct
		}

		if promotion.UUID == "" {
			promotion = rawPromotion
		}
	}

	for _, product := range mapProduct {
		if variants, ok := mapProductIdVariants[product.UUID]; ok {
			product.ProductVariants = variants
			products = append(products, product)
		}
	}

	d.RedisInfra.Client.Set(promotionString, promotion, utils.DefaultRedisTimeOut)

	span.End()

	promotion.Products = products
	res.CorePromotion = promotion
	return res, traceID, http.StatusOK, nil
}

func (d *DiscountService) UpdateDiscount(ctx context.Context, id string, data model.DiscountRequest) (traceID string, statusCode int, err error) {
	ctx, span := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount")
	traceID = span.SpanContext().TraceID().String()

	startTime := data.StartTime.Format(time.RFC3339)
	endTime := data.EndTime.Format(time.RFC3339)

	promotionInterface := []interface{}{
		data.Name, data.Code, startTime, endTime, data.UsageQuantity, data.UsageLimitPerUser, id,
	}

	promotionSelectString := utils.PrepareSelectQueryForUpdate(utils.PromotionTableName, []string{"1"})

	promotionString := utils.PrepareUpdateQuery(utils.PromotionTableName, utils.PromotionColumnsListForUpdate)

	_, updateTxSpan := d.OtelInfra.Tracer.Start(ctx, "UpdateDiscount-Tx")

	tx, err := d.PostgresInfra.DbWritePool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}
	defer tx.Rollback(ctx)

	updateTxSpan.End()

	_, dbSpan := d.OtelInfra.Tracer.Start(ctx, "UpdateDiscount-DB")
	dbSpan.AddEvent("update promotion", trace.WithAttributes(
		attribute.String("query", promotionString),
		attribute.String("params", fmt.Sprintf("%v", promotionInterface)),
	))

	_, err = tx.Exec(ctx, promotionSelectString, id)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	_, err = tx.Exec(ctx, promotionString, promotionInterface...)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	_, err = tx.Exec(ctx, "DELETE FROM products WHERE promotion_id = $1", id)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	_, err = tx.Exec(ctx, "DELETE FROM product_variants WHERE promotion_id = $1", id)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	// products
	var productsInterface [][]interface{}
	var variantsInterface [][]interface{}

	for i := 0; i < len(data.Products); i++ {
		data.Products[i].UUID = uuid.New().String()
		data.Products[i].PromotionID = id
		productInterface := []interface{}{
			data.Products[i].UUID,
			data.Products[i].PromotionID,
			data.Products[i].SKU,
			data.Products[i].Name,
			data.Products[i].PurchaseLimit,
		}
		productsInterface = append(productsInterface, productInterface)

		for j := 0; j < len(data.Products[i].ProductVariants); j++ {
			data.Products[i].ProductVariants[j].UUID = uuid.New().String()
			data.Products[i].ProductVariants[j].PromotionID = id
			data.Products[i].ProductVariants[j].ProductID = data.Products[i].UUID
			variantInterface := []interface{}{
				data.Products[i].ProductVariants[j].UUID,
				data.Products[i].ProductVariants[j].PromotionID,
				data.Products[i].ProductVariants[j].ProductID,
				data.Products[i].ProductVariants[j].SKU,
				data.Products[i].ProductVariants[j].Name,
				data.Products[i].ProductVariants[j].DiscountedPrice,
				data.Products[i].ProductVariants[j].DiscountedPercentage,
				data.Products[i].ProductVariants[j].StockLimit,
				data.Products[i].ProductVariants[j].IsActive,
			}
			variantsInterface = append(variantsInterface, variantInterface)
		}
	}

	productString := utils.PrepareInsertQuery(utils.ProductTableName, utils.ProductColumnsListForInsert, productsInterface)

	var productParams []interface{}
	for _, product := range productsInterface {
		productParams = append(productParams, product...)
	}

	dbSpan.AddEvent("insert products", trace.WithAttributes(
		attribute.String("query", productString),
		attribute.String("params", fmt.Sprintf("%v", productParams)),
	))

	_, err = tx.Exec(ctx, productString, productParams...)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	// variants
	variantString := utils.PrepareInsertQuery(utils.VariantTableName, utils.VariantColumnsListForInsert, variantsInterface)

	var variantParams []interface{}
	for _, variant := range variantsInterface {
		variantParams = append(variantParams, variant...)
	}

	dbSpan.AddEvent("insert variants", trace.WithAttributes(
		attribute.String("query", variantString),
		attribute.String("params", fmt.Sprintf("%v", variantParams)),
	))

	_, err = tx.Exec(ctx, variantString, variantParams...)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	dbSpan.End()
	span.End()

	return traceID, http.StatusOK, nil
}

func (d *DiscountService) ListDiscounts(ctx context.Context, query core_model.CoreQuery) (res []model.DiscountResponse, pagination core_model.Pagination, traceID string, statusCode int, err error) {
	ctx, span := d.OtelInfra.Tracer.Start(ctx, "ListDiscounts")
	traceID = span.SpanContext().TraceID().String()

	pagination.Page = query.Page
	countString := utils.PrepareSelectCountQuery(utils.PromotionTableName)

	intVariable := 0
	promotionInterface := []interface{}{}
	countInterface := []interface{}{}
	promotionString := utils.PrepareSelectQuery(utils.PromotionTableName, utils.PromotionColumnsListForSelect)

	if query.ShopID != "" {
		intVariable++
		promotionString += fmt.Sprintf("AND shop_id = $%d ", intVariable)
		promotionInterface = append(promotionInterface, query.ShopID)

		countString += fmt.Sprintf("AND shop_id = $%d ", intVariable)
		countInterface = append(countInterface, query.ShopID)
	}

	if query.Cursor != "" && query.Page != "" {
		intVariable++
		promotionString += fmt.Sprintf("AND id > $%d ", intVariable)
		promotionInterface = append(promotionInterface, query.Cursor)
	}

	if query.Sort == "" {
		query.Sort = utils.DefaultOrder
	}
	promotionString += fmt.Sprintf(`ORDER BY %s `, query.Sort)

	intVariable++
	promotionString += fmt.Sprintf(`LIMIT $%d `, intVariable)
	if query.Size == "" {
		query.Size = utils.DefaultSize
	}
	promotionInterface = append(promotionInterface, query.Size)
	pagination.Size = query.Size

	intVariable++
	promotionString += fmt.Sprintf(`OFFSET $%d `, intVariable)
	if query.Page == "" {
		query.Page = utils.DefaultPage
	}
	offset := (utils.StringToInt(query.Page) - 1) * utils.StringToInt(query.Size)
	promotionInterface = append(promotionInterface, utils.IntToString(offset))
	pagination.Page = query.Page

	cache, err := d.RedisInfra.Client.Get(promotionString).Bytes()
	if err == nil {
		err = json.Unmarshal(cache, &res)
		if err == nil {
			return res, pagination, traceID, http.StatusOK, nil
		}
	}

	// count
	err = d.PostgresInfra.DbReadPool.QueryRow(ctx, countString, countInterface...).Scan(&pagination.TotalItems)
	if err != nil {
		return res, pagination, traceID, http.StatusInternalServerError, err
	}

	_, dbSpan := d.OtelInfra.Tracer.Start(ctx, "ListDiscounts-DB")
	dbSpan.AddEvent("get promotions", trace.WithAttributes(
		attribute.String("query", promotionString),
		attribute.String("params", fmt.Sprintf("%v", promotionInterface)),
	))

	promotionRows, err := d.PostgresInfra.DbReadPool.Query(ctx, promotionString, promotionInterface...)
	if err != nil {
		return res, pagination, traceID, http.StatusInternalServerError, err
	}
	defer promotionRows.Close()

	fmt.Printf("SQL: %s\nParams: %#v\n", promotionString, promotionInterface)

	promotionsId := []string{}
	promotions := []model.DiscountResponse{}
	for promotionRows.Next() {
		var rawPromotion core_model.CorePromotion
		err = promotionRows.Scan(
			&rawPromotion.ID,
			&rawPromotion.CreatedAt,
			&rawPromotion.UpdatedAt,
			&rawPromotion.DeletedAt,
			&rawPromotion.UUID,
			&rawPromotion.Name,
			&rawPromotion.PromotionType,
			&rawPromotion.Code,
			&rawPromotion.StartTime,
			&rawPromotion.EndTime,
			&rawPromotion.ShopID,
			&rawPromotion.UsageQuantity,
			&rawPromotion.UsageLimitPerUser,
		)

		if err != nil {
			return res, pagination, traceID, http.StatusInternalServerError, err
		}

		promotionsId = append(promotionsId, fmt.Sprintf(`'%s'`, rawPromotion.UUID))

		promotions = append(promotions, model.DiscountResponse{
			CorePromotion: rawPromotion,
		})
	}

	if len(promotionsId) == 0 {
		dbSpan.End()
		return promotions, pagination, traceID, http.StatusOK, nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error)

	mapPromotionIdProducts := make(map[string]core_model.CoreProduct)
	mapProductIdVariants := make(map[string][]core_model.CoreProductVariant)

	go func(promotionsId []string) {
		defer wg.Done()

		productString := utils.PrepareSelectQuery(utils.ProductTableName, utils.ProductColumnsListForSelect)
		productString += fmt.Sprintf("AND promotion_id IN (%s) ", strings.Join(promotionsId, ", "))

		dbSpan.AddEvent("get products", trace.WithAttributes(
			attribute.String("params", productString),
		))

		productRows, err := d.PostgresInfra.DbReadPool.Query(ctx, productString)
		if err != nil {
			errChan <- err
			return
			// return res, http.StatusInternalServerError, err
		}
		defer productRows.Close()

		for productRows.Next() {
			var rawProduct core_model.CoreProduct
			err = productRows.Scan(
				&rawProduct.ID,
				&rawProduct.UUID,
				&rawProduct.PromotionID,
				&rawProduct.SKU,
				&rawProduct.Name,
				&rawProduct.PurchaseLimit,
			)
			if err != nil {
				errChan <- err
				return
				// return res, http.StatusInternalServerError, err
			}

			key := rawProduct.PromotionID + "_" + rawProduct.UUID
			if _, ok := mapPromotionIdProducts[key]; !ok {
				mapPromotionIdProducts[key] = rawProduct
			}

		}
	}(promotionsId)

	go func(promotionsId []string) {
		defer wg.Done()

		variantString := utils.PrepareSelectQuery(utils.VariantTableName, utils.VariantColumnsListForSelect)

		variantString += fmt.Sprintf("AND promotion_id IN (%s) ", strings.Join(promotionsId, ", "))

		dbSpan.AddEvent("get variants", trace.WithAttributes(
			attribute.String("query", variantString),
		))
		variantRows, err := d.PostgresInfra.DbReadPool.Query(ctx, variantString)
		if err != nil {
			errChan <- err
			return
			// return res, http.StatusInternalServerError, err
		}
		defer variantRows.Close()

		for variantRows.Next() {
			var rawProductVariant core_model.CoreProductVariant
			err = variantRows.Scan(
				&rawProductVariant.ID,
				&rawProductVariant.PromotionID,
				&rawProductVariant.ProductID,
				&rawProductVariant.SKU,
				&rawProductVariant.Name,
				&rawProductVariant.DiscountedPrice,
				&rawProductVariant.DiscountedPercentage,
				&rawProductVariant.StockLimit,
				&rawProductVariant.IsActive,
			)
			if err != nil {
				errChan <- err
				return
			}

			if _, ok := mapProductIdVariants[rawProductVariant.ProductID]; !ok {
				mapProductIdVariants[rawProductVariant.ProductID] = []core_model.CoreProductVariant{}
			}

			mapProductIdVariants[rawProductVariant.ProductID] = append(mapProductIdVariants[rawProductVariant.ProductID], rawProductVariant)
		}
	}(promotionsId)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return res, pagination, traceID, http.StatusInternalServerError, err
		}
	}

	dbSpan.End()
	_, mappingSpan := d.OtelInfra.Tracer.Start(ctx, "ListDiscounts-Mapping")
	for i := 0; i < len(promotions); i++ {
		var products []core_model.CoreProduct

		for key, product := range mapPromotionIdProducts {
			parts := strings.Split(key, "_")
			if parts[0] != promotions[i].UUID {
				continue
			} else {
				if variant, ok := mapProductIdVariants[parts[1]]; ok {
					product.ProductVariants = variant
					products = append(products, product)
				}
			}

		}

		promotions[i].Products = products
	}

	mappingSpan.End()

	d.RedisInfra.Client.Set(promotionString, promotions, utils.DefaultRedisTimeOut)

	span.End()
	return promotions, pagination, traceID, http.StatusOK, nil
}

func (d *DiscountService) DeleteDiscount(ctx context.Context, id string) (traceID string, statusCode int, err error) {
	ctx, span := d.OtelInfra.Tracer.Start(ctx, "CreateDiscount")
	traceID = span.SpanContext().TraceID().String()

	_, dbSpan := d.OtelInfra.Tracer.Start(ctx, "DeleteDiscount-DB")

	queryString := `UPDATE promotions SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`
	_, err = d.PostgresInfra.DbWritePool.Exec(context.Background(), queryString, id)
	if err != nil {
		return traceID, http.StatusInternalServerError, err
	}

	dbSpan.End()
	span.End()

	return traceID, http.StatusOK, nil
}
