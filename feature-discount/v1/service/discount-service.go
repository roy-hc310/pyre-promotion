package service

import (
	"encoding/json"
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
)

type DiscountService struct {
	PostgresInfra *infrastructure.PostgresInfra
	RedisInfra    *infrastructure.RedisInfra
}

func NewDiscountService(postgresInfra *infrastructure.PostgresInfra, redisInfra *infrastructure.RedisInfra) *DiscountService {
	return &DiscountService{
		PostgresInfra: postgresInfra,
		RedisInfra:    redisInfra,
	}
}

func (d *DiscountService) CreateDiscount(data model.DiscountRequest) (res core_model.CoreIdResponse, statusCode int, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultContextTimeOut)
	defer cancel()

	tx, err := d.PostgresInfra.DbWritePool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return res, http.StatusInternalServerError, err
	}
	defer tx.Rollback(ctx)

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

	_, err = tx.Exec(ctx, promotionString, promotionParams...)
	if err != nil {
		return res, http.StatusInternalServerError, err
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

	_, err = tx.Exec(ctx, productString, productParams...)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	// variants
	variantString := utils.PrepareInsertQuery(utils.VariantTableName, utils.VariantColumnsListForInsert, variantsInterface)

	var variantParams []interface{}
	for _, variant := range variantsInterface {
		variantParams = append(variantParams, variant...)
	}

	_, err = tx.Exec(ctx, variantString, variantParams...)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res.Id = promotionUUID

	return res, http.StatusOK, nil
}

func (d *DiscountService) CreateBulkDiscount(data []model.DiscountRequest) (res []core_model.CoreIdResponse, statusCode int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultContextTimeOut)
	defer cancel()

	tx, err := d.PostgresInfra.DbWritePool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return res, http.StatusInternalServerError, err
	}
	defer tx.Rollback(ctx)

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
	_, err = tx.Exec(ctx, promotionString, promotionParams...)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	// products
	productString := utils.PrepareInsertQuery(utils.ProductTableName, utils.ProductColumnsListForInsert, productsInterface)
	var productParams []interface{}
	for _, product := range productsInterface {
		productParams = append(productParams, product...)
	}
	_, err = tx.Exec(ctx, productString, productParams...)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	// variants
	variantString := utils.PrepareInsertQuery(utils.VariantTableName, utils.VariantColumnsListForInsert, variantsInterface)
	var variantParams []interface{}
	for _, variant := range variantsInterface {
		variantParams = append(variantParams, variant...)
	}
	_, err = tx.Exec(ctx, variantString, variantParams...)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}

	res = response

	return res, http.StatusOK, nil
}

func (d *DiscountService) DetailDiscount(id string) (res model.DiscountResponse, statusCode int, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultContextTimeOut)
	defer cancel()

	promotionString := fmt.Sprintf(`SELECT p.id, p.created_at, p.updated_at, p.deleted_at, p.uuid, p.name, p.promotion_type, p.code, p.start_time, p.end_time, p.shop_id, p.usage_quantity, p.usage_limit_per_user, pr.id, pr.uuid, pv.promotion_id, pr.sku, pr.name, pr.purchase_limit,pv.id, pv.uuid, pv.promotion_id, pv.product_id, pv.sku, pv.name, pv.discounted_price, pv.discounted_percentage, pv.stock_limit, pv.is_active FROM promotion.promotions p LEFT JOIN promotion.products pr on pr.promotion_id = p.uuid LEFT JOIN promotion.product_variants pv on pv.product_id = pr.uuid WHERE p.uuid = '%s' AND p.deleted_at IS NULL;`, id)

	cache, err := d.RedisInfra.Client.Get(promotionString).Bytes()
	if err == nil {
		err = json.Unmarshal(cache, &res)
		if err == nil {
			return res, http.StatusOK, nil
		}
	}

	rows, err := d.PostgresInfra.DbReadPool.Query(ctx, promotionString)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}
	defer rows.Close()

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
			return res, http.StatusInternalServerError, err
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

	promotion.Products = products
	res.CorePromotion = promotion
	return res, http.StatusOK, nil
}

func (d *DiscountService) UpdateDiscount(id string, data model.DiscountRequest) (statusCode int, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultContextTimeOut)
	defer cancel()

	tx, err := d.PostgresInfra.DbWritePool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer tx.Rollback(ctx)

	startTime := data.StartTime.Format(time.RFC3339)
	endTime := data.EndTime.Format(time.RFC3339)

	promotionInterface := []interface{}{
		data.Name, data.Code, startTime, endTime, data.UsageQuantity, data.UsageLimitPerUser, id,
	}

	promotionString := utils.PrepareUpdateQuery(utils.PromotionTableName, utils.PromotionColumnsListForUpdate)

	_, err = tx.Exec(ctx, promotionString, promotionInterface...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = tx.Exec(ctx, "DELETE FROM products WHERE promotion_id = $1", id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = tx.Exec(ctx, "DELETE FROM product_variants WHERE promotion_id = $1", id)
	if err != nil {
		return http.StatusInternalServerError, err
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

	_, err = tx.Exec(ctx, productString, productParams...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// variants
	variantString := utils.PrepareInsertQuery(utils.VariantTableName, utils.VariantColumnsListForInsert, variantsInterface)

	var variantParams []interface{}
	for _, variant := range variantsInterface {
		variantParams = append(variantParams, variant...)
	}

	_, err = tx.Exec(ctx, variantString, variantParams...)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (d *DiscountService) ListDiscounts(query core_model.CoreQuery) (res []model.DiscountResponse, statusCode int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultContextTimeOut)
	defer cancel()

	// promotionString := fmt.Sprintf(`SELECT id, created_at, updated_at, deleted_at, uuid, name, promotion_type, code, start_time, end_time, shop_id, usage_quantity, usage_limit_per_user FROM promotions WHERE deleted_at IS NULL ORDER BY %s LIMIT %s`, query.Sort, query.Size)

	promotionString := utils.PrepareSelectQuery(utils.PromotionTableName, utils.PromotionColumnsListForSelect)

	if query.ShopID != "" {
		promotionString += fmt.Sprintf("AND shop_id = '%s' ", query.ShopID)
	}

	if query.Cursor != "" {
		promotionString += fmt.Sprintf(" AND id > '%s' ", query.Cursor)
	}

	if query.Sort != "" {
		promotionString += fmt.Sprintf(`ORDER BY %s `, query.Sort)
	} else {
		promotionString += `ORDER BY id DESC `
	}

	if query.Size != "" {
		promotionString += fmt.Sprintf(`LIMIT %s;`, query.Size)
	} else {
		promotionString += fmt.Sprintf(`LIMIT %s;`, utils.DefaultSize)
	}

	cache, err := d.RedisInfra.Client.Get(promotionString).Bytes()
	if err == nil {
		err = json.Unmarshal(cache, &res)
		if err == nil {
			return res, http.StatusOK, nil
		}
	}

	promotionRows, err := d.PostgresInfra.DbReadPool.Query(ctx, promotionString)
	if err != nil {
		return res, http.StatusInternalServerError, err
	}
	defer promotionRows.Close()

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
			return res, http.StatusInternalServerError, err
		}

		promotionsId = append(promotionsId, fmt.Sprintf(`'%s'`, rawPromotion.UUID))

		promotions = append(promotions, model.DiscountResponse{
			CorePromotion: rawPromotion,
		})
	}

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error)

	mapPromotionIdProducts := make(map[string]core_model.CoreProduct)
	mapProductIdVariants := make(map[string][]core_model.CoreProductVariant)

	
	go func() {
		defer wg.Done()

		productString := utils.PrepareSelectQuery(utils.ProductTableName, utils.ProductColumnsListForSelect)
		productString += fmt.Sprintf("AND promotion_id IN (%s) ", strings.Join(promotionsId, ", "))

		productRows, err := d.PostgresInfra.DbReadPool.Query(ctx, productString)
		if err != nil {
			errChan <- err
			return
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
			}

			key := rawProduct.PromotionID + "_" + rawProduct.UUID
			if _, ok := mapPromotionIdProducts[key]; !ok {
				mapPromotionIdProducts[key] = rawProduct
			}

		}
	}()

	go func() {
		defer wg.Done()

		variantString := utils.PrepareSelectQuery(utils.VariantTableName, utils.VariantColumnsListForSelect)

		variantString += fmt.Sprintf("AND promotion_id IN (%s) ", strings.Join(promotionsId, ", "))

		variantRows, err := d.PostgresInfra.DbReadPool.Query(ctx, variantString)
		if err != nil {
			errChan <- err
			return
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
	}()

	// go func() {
		wg.Wait()
		close(errChan)
	// }()

	for err := range errChan {
		if err != nil {
			return res, http.StatusInternalServerError, err
		}
	}

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

	d.RedisInfra.Client.Set(promotionString, promotions, utils.DefaultRedisTimeOut)

	return promotions, http.StatusOK, nil
}

func (d *DiscountService) DeleteDiscount(id string) (statusCode int, err error) {

	queryString := fmt.Sprintf(`UPDATE promotions SET deleted_at = now() WHERE id = '%s' AND deleted_at IS NULL`, id)
	_, err = d.PostgresInfra.DbWritePool.Exec(context.Background(), queryString)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
