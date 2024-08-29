package utils

import "time"

const (
	XPlatformKey = "x-platform-key"
	XShopId      = "x-shop-id"
)

const (
	DefaultPage = 1
	DefaultSize = "10"
)

const (
	DefaultContextTimeOut = time.Second * 10
	DefaultRedisTimeOut   = time.Second * 10
)

const (
	TopicCreateBulkDiscount = "topic-create-bulk-discount"
)

const (
	PromotionTableName = "promotions"
	ProductTableName   = "products"
	VariantTableName   = "product_variants"
)

const (
	PromotionTypeDiscount        = "discount"
	PromotionTypeFlashSale       = "flashsale"
	PromotionTypeAddonDeal       = "addon-deal"
	PromotionTypeBundleDeal      = "bundle-deal"
	PromotionTypeVoucher         = "voucher"
	PromotionTypeVoucherShipping = "voucher-shipping"
)

var PromotionColumnsListForInsert = []string{"uuid", "name", "promotion_type", "code", "start_time", "end_time", "shop_id", "usage_quantity", "usage_limit_per_user"}
var ProductColumnsListForInsert = []string{"uuid", "promotion_id", "sku", "name", "purchase_limit"}
var VariantColumnsListForInsert = []string{"uuid", "promotion_id", "product_id", "sku", "name", "discounted_price", "discounted_percentage", "stock_limit", "is_active"}

var PromotionColumnsListForUpdate = []string{"name", "code", "start_time", "end_time", "usage_quantity", "usage_limit_per_user"}

var PromotionColumnsListForSelect = []string{"id", "created_at", "updated_at", "deleted_at", "uuid", "name", "promotion_type", "code", "start_time", "end_time", "shop_id", "usage_quantity", "usage_limit_per_user"}
var ProductColumnsListForSelect = []string{"id", "uuid", "promotion_id", "sku", "name", "purchase_limit"}
var VariantColumnsListForSelect = []string{"id", "promotion_id", "product_id", "sku", "name", "discounted_price", "discounted_percentage", "stock_limit", "is_active"}
