package web

type Sort string

const (
	ASC Sort = "ASC"
	DESC Sort = "DESC"
)

type ProductGetRequest struct{
	Page *int `form:"page" binding:"omitempty,gt=0"`
	Sort *string `form:"sort" binding:"omitempty,oneof=ASC DESC"`
	Limit *int `form:"limit" binding:"omitempty,oneof=10 50 100,gte=0"`
	Order *string `form:"order" binding:"omitempty"`
	Search *string `form:"search" binding:"omitempty"`
}