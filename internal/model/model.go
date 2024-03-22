package model

type AccountCreateRequest struct {
	AccountID int      `json:"account_id" validate:"required"`
	Limit     int      `json:"limit" validate:"required"`
	Products  []string `json:"products" validate:"required"`
}

type AccountListRequest struct {
	Product string     `query:"product"`
	OrderBy OrderField `query:"order_by"`
	Limit   int        `query:"limit"`
	Page    int        `query:"page"`
}

type OrderField struct {
	AccountID string `query:"account_id"`
}

type AccountUpdateRequest struct {
	Limit    int      `json:"limit" validate:"required"`
	Products []string `json:"products" validate:"required"`
}

type AccountResponse struct {
	AccountID int      `json:"account_id"`
	Limit     int      `json:"limit"`
	Products  []string `json:"products"`
}
