package repository

import (
	"belajar-go-restful-api/helper"
	models "belajar-go-restful-api/model/domain"
	"belajar-go-restful-api/model/web"
	"fmt"
	"net/http"
	"strings"
)

func (r *ProductRepositoryImpl) FindAll(params web.ProductGetRequest) ([]models.Product, *models.Metadata, error){
	var products []models.Product;

	trDB := r.db.Model(&products)
	
	if params.Page == nil{
		defaultPage := 1
		params.Page = &defaultPage
	}
	
	
	if params.Limit == nil{
		defaultLimit := 10
		params.Limit = &defaultLimit
	}
	

	if params.Sort == nil{
		defaultSort := "DESC"
		params.Sort = &defaultSort
	}

	if params.Order == nil{
		defaultOrder := "created_at"
		params.Order = &defaultOrder
	}

	if params.Search != nil{
		src := strings.Join([]string{"%", *params.Search, "%"}, "")
		trDB = trDB.Where("product_name ILIKE ?", src)
	}

	trDB = trDB.Order(fmt.Sprintf("%s %s", *params.Order, *params.Sort))

	totalProduct := int64(0)
	err := trDB.Count(&totalProduct).Error

	if err != nil{
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, nil, err
	}

	offset := (*params.Page - 1) * *params.Limit

	trDB = trDB.Limit(*params.Limit).Offset(offset)

	pagination := helper.PaginationData(int64(*params.Page), float64(*params.Limit), float64(totalProduct));

	err = trDB.Find(&products).Error;

	if err != nil{
		return nil, nil, err
	}

	return products, pagination, nil
}

func (r *ProductRepositoryImpl) FindByID(id int) (*models.Product, error){
	var product *models.Product

	err := r.db.First(&product, id).Error

	if err != nil{
		err = helper.SetError(http.StatusNotFound, "Data product not found!")
		return nil, err
	}

	return product, nil
}

func (r *ProductRepositoryImpl) Create(productRequest *models.Product) (*models.Product, error){
	var product *models.Product
	err := r.db.Where("product_name = ? AND description = ? AND price = ?", productRequest.ProductName, productRequest.Description, productRequest.Price).First(&product).Error;

	if err == nil{
		if product != nil{
			err := helper.SetError(http.StatusConflict, "Data product already exists!")
			return nil, err
		}
		err := helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	err = r.db.Create(&productRequest).Error;

	if err != nil{
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	return productRequest, nil
}

func (r *ProductRepositoryImpl) Update(productUpdateRequest *models.Product) (*models.Product, error){
	var product *models.Product
	err := r.db.Where("product_name = ? AND description = ? AND price = ? AND id != ?", productUpdateRequest.ProductName, productUpdateRequest.Description, productUpdateRequest.Price, productUpdateRequest.ID).First(&product).Error;

	if err == nil{
		if product != nil{
			err := helper.SetError(http.StatusConflict, "Data product already exists")
			return nil, err
		}
		err := helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	err = r.db.Updates(&productUpdateRequest).Error

	if err != nil{
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	return productUpdateRequest, nil
}

func (r *ProductRepositoryImpl) Delete(id int) error{
	err := r.db.Delete(&models.Product{ID: uint(id)}).Error

	if err != nil{
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return err
	}

	return nil
}