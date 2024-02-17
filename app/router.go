package app

import (
	handler "belajar-go-restful-api/handler/product"
	"belajar-go-restful-api/middleware"

	"github.com/gin-gonic/gin"
)

func ProductRouter(productHandler handler.ProductHandlerImpl, app *gin.Engine) *gin.RouterGroup{
	authorized := app.Group("/api/v1", middleware.AuthMiddleware())
	authorized.GET("/products/:id", productHandler.FindByID)
	authorized.GET("/storage/:filename", productHandler.GetProductImage)
	authorized.POST("/products", productHandler.Create)
	authorized.PATCH("/products/:id", productHandler.Update)
	authorized.DELETE("/products/:id", productHandler.Delete)

	unauthorized := app.Group("/api/v1")
	unauthorized.GET("/products",productHandler.FindAll)

	return unauthorized
}