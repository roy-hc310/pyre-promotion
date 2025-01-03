package model

import "github.com/jackc/pgx/v5/pgtype"

type CoreModel struct {
	ID        int              `json:"id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
	UUID      string           `json:"uuid"`
	TraceID   string           `json:"trace_id"`
}

type CoreQuery struct {
	Page string `form:"page" json:"page"`
	Size string `form:"size" json:"size"`

	Sort   string `form:"sort" json:"sort"`
	Search string `form:"search" json:"search"`
	Cursor string `form:"cursor" json:"cursor"`
	ShopID string `form:"shop_id" json:"shop_id"`
}

type CoreResponseObject struct {
	Data      interface{} `json:"data"`
	TraceID   string      `json:"trace_id"`
	Succeeded bool        `json:"succeeded"`
	Errors    []string    `json:"errors"`
}

type CoreResponseArray struct {
	CoreResponseObject
	TotalItems int `json:"total_items"`
}
