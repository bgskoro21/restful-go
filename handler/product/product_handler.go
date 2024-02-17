package handler

import (
	service "belajar-go-restful-api/service/product"
)

type ProductHandlerImpl struct{
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandlerImpl{
	return &ProductHandlerImpl{productService} 
}