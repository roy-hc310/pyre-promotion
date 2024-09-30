package test

import (
	"context"
	"net/http"
	"pyre-promotion/core-internal/infrastructure"
	"pyre-promotion/feature-discount/v1/service"
	"testing"
	"time"

	core_model "pyre-promotion/core-internal/model"
	"pyre-promotion/feature-discount/v1/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

type MockPostgresInfra struct {
	mock.Mock
	infrastructure.PostgresInfra
}

func NewMockPostgresInfra() *MockPostgresInfra {
	mockPostgresInfra := &MockPostgresInfra{}
	mockTx := new(MockTx)

	mockPostgresInfra.On("DbWritePool").Return(mockTx)
	mockTx.On("BeginTx", mock.Anything, mock.Anything).Return(mockTx, nil)
	return mockPostgresInfra
}

func (m *MockPostgresInfra) BeginTx(ctx context.Context) (*pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(*pgx.Tx), args.Error(1)
}

type MockRedisInfra struct {
	mock.Mock
	infrastructure.RedisInfra
}


func NewMockRedisInfra() *MockRedisInfra {
	return &MockRedisInfra{}
}

func SetupDiscountTest() (*service.DiscountService, *MockPostgresInfra, *MockRedisInfra) {
	mockPostgres := NewMockPostgresInfra()
	mockRedis := NewMockRedisInfra()
	discountService := service.NewDiscountService(&mockPostgres.PostgresInfra, &mockRedis.RedisInfra)

	return discountService, mockPostgres, mockRedis
}

func TestCreateDiscount(t *testing.T) {
	discountService, _, _ := SetupDiscountTest()

	discountReqs := model.DiscountRequest{
		CorePromotion: core_model.CorePromotion{
			CoreModel: core_model.CoreModel{
				ID:        1,
				CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
				DeletedAt: pgtype.Timestamp{Valid: false},
				UUID:      "promotion-uuid",
			},
			Name:              "Summer Sale",
			PromotionType:     "percentage_discount",
			Code:              "SUMMER2024",
			StartTime:         time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC),
			EndTime:           time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC),
			ShopID:            "shop-uuid",
			UsageQuantity:     1000,
			UsageLimitPerUser: 5,
			Products: []core_model.CoreProduct{
				{
					CoreModel: core_model.CoreModel{
						ID:        101,
						CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
						UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
						DeletedAt: pgtype.Timestamp{Valid: false},
						UUID:      "product-uuid-1",
					},
					PromotionID:   "promotion-uuid",
					SKU:           "PROD001",
					Name:          "Product 1",
					PurchaseLimit: 10,
					ProductVariants: []core_model.CoreProductVariant{
						{
							CoreModel: core_model.CoreModel{
								ID:        1001,
								CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
								UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
								DeletedAt: pgtype.Timestamp{Valid: false},
								UUID:      "variant-uuid-1",
							},
							PromotionID:          "promotion-uuid",
							ProductID:            "product-uuid-1",
							SKU:                  "PROD001-V1",
							Name:                 "Product 1 - Variant 1",
							DiscountedPrice:      50.0,
							DiscountedPercentage: 10.0,
							StockLimit:           100,
							IsActive:             true,
						},
					},
				},
				{
					CoreModel: core_model.CoreModel{
						ID:        102,
						CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
						UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
						DeletedAt: pgtype.Timestamp{Valid: false},
						UUID:      "product-uuid-2",
					},
					PromotionID:   "promotion-uuid",
					SKU:           "PROD002",
					Name:          "Product 2",
					PurchaseLimit: 5,
					ProductVariants: []core_model.CoreProductVariant{
						{
							CoreModel: core_model.CoreModel{
								ID:        1002,
								CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
								UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
								DeletedAt: pgtype.Timestamp{Valid: false},
								UUID:      "variant-uuid-2",
							},
							PromotionID:          "promotion-uuid",
							ProductID:            "product-uuid-2",
							SKU:                  "PROD002-V1",
							Name:                 "Product 2 - Variant 1",
							DiscountedPrice:      70.0,
							DiscountedPercentage: 15.0,
							StockLimit:           50,
							IsActive:             true,
						},
					},
				},
			},
		},
	}

	_, statusCode, err := discountService.CreateDiscount(discountReqs)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
}