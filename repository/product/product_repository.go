package repository

import (
	models "belajar-go-restful-api/model/domain"
	"belajar-go-restful-api/model/web"

	"gorm.io/gorm"
)

type ProductRepository interface{
	FindAll(params web.ProductGetRequest) ([]models.Product, *models.Metadata, error)
	FindByID(id int) (*models.Product, error)
	Create(productRequest *models.Product) (*models.Product, error)
	Update(productUpdateRequest *models.Product) (*models.Product, error)
	Delete(id int) error
}

type ProductRepositoryImpl struct{
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl{
	return &ProductRepositoryImpl{db}
}