package web

import "mime/multipart"

type ProductUpdateRequest struct{
	ID int `form:"id" binding:"omitempty"`
	ProductName string `form:"product_name" binding:"required,min=1,max=100"`
	Description string `form:"description" binding:"required"` 
	Price int `form:"price_product" binding:"required"`
	ProductImage *multipart.FileHeader `form:"product_image" binding:"omitempty"`
}