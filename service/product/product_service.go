package service

import (
	models "belajar-go-restful-api/model/domain"
	"belajar-go-restful-api/model/web"
	repository "belajar-go-restful-api/repository/product"
)

type ProductService interface{
	FindAll(params web.ProductGetRequest) ([]models.Product, *models.Metadata, error)
	FindByID(id int) (*models.Product, error)
	Create(productRequest web.ProductCreateRequest) (*models.Product, error)
	Update(productRequest web.ProductUpdateRequest) (*models.Product, error)
	Delete(id int) error
}

type ProductServiceImpl struct{
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) *ProductServiceImpl{
	return &ProductServiceImpl{productRepository}
}