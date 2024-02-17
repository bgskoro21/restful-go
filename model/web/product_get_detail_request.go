package web

type ProductGetDetailRequest struct{
	ID int `form:"id" binding:"required,gt:0"`
}