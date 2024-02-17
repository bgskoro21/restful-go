package main

import (
	"belajar-go-restful-api/app"
	handler "belajar-go-restful-api/handler/product"
	repository "belajar-go-restful-api/repository/product"
	service "belajar-go-restful-api/service/product"

	"github.com/gin-gonic/gin"
)

func main(){
	db := app.NewDB()

	r := gin.Default()

	productRepository := repository.NewProductRepository(db)

	productService := service.NewProductService(productRepository)

	productHandler := handler.NewProductHandler(productService)

	app.ProductRouter(*productHandler, r)

	r.Run()
}