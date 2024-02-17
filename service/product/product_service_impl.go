package service

import (
	"belajar-go-restful-api/helper"
	models "belajar-go-restful-api/model/domain"
	"belajar-go-restful-api/model/web"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/drive/v3"
)

func (p *ProductServiceImpl) FindAll(params web.ProductGetRequest) ([]models.Product, *models.Metadata, error){
	products, metadata, err := p.productRepository.FindAll(params)

	if err != nil{
		return nil, nil, err
	}

	return products, metadata, nil
}

func (p *ProductServiceImpl) FindByID(id int) (*models.Product, error){
	product, err := p.productRepository.FindByID(id)

	if err != nil{
		return nil, err
	}

	return product, err
}

func (p *ProductServiceImpl) Create(productRequest web.ProductCreateRequest) (*models.Product, error){
	file := productRequest.ProductImage

	filename := strings.ReplaceAll(file.Filename, " ", "_")

	src,err := file.Open()

	if err != nil{
		return nil, err
	}

	defer src.Close()

	srv, err := helper.GetDriveService()

	if err != nil{
		fmt.Println(err.Error())
	}

	folderId := "1pJsJk8UO0bB3tne-doVBLE5I0i6p0C8r"

	// Step 4: create the file and upload
	files, err := helper.CreateFile(srv, filename, "application/octet-stream", src, folderId)

	if err != nil {
		fmt.Println("Error creating file:", err)
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	err = helper.SetFilePermissions(srv, files.Id)

	if err != nil{
		fmt.Println("Error setting file permissions:", err)
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	url := fmt.Sprintf("https://drive.google.com/uc?id=%s", files.Id)

	request := &models.Product{
		ProductName: productRequest.ProductName,
		Description: productRequest.Description,
		Price: productRequest.Price,
		ProductImage: &url,
	}

	product,err := p.productRepository.Create(request)

	if err != nil{
		errDelete := helper.DeleteFile(srv, files.Id)
		if errDelete != nil {
			fmt.Println("Error deleting file:", errDelete)
			return nil, errDelete
		}
		return nil, err
	}

	return product, nil
} 

func (p *ProductServiceImpl) Update(productRequest web.ProductUpdateRequest) (*models.Product, error){
	var (
		files *drive.File
		fileID string
	)

	product, err := p.FindByID(productRequest.ID)

	if err != nil{
		return nil, err
	}

	srv, err := helper.GetDriveService()

	if err != nil{
		fmt.Println(err.Error())
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	if productRequest.ProductImage != nil{
		file := productRequest.ProductImage

		filename := strings.ReplaceAll(file.Filename, " ", "_")

		src,err := file.Open()

		if err != nil{
			return nil, err
		}

		defer src.Close()

	    fileID, err = helper.GetFileIDFromURL(*product.ProductImage)

		if err != nil{
			err = helper.SetError(http.StatusInternalServerError, err.Error())
			return nil, err
		}

		folderId := "1pJsJk8UO0bB3tne-doVBLE5I0i6p0C8r"

		// Step 4: create the file and upload
		files, err = helper.CreateFile(srv, filename, "application/octet-stream", src, folderId)

		if err != nil {
			err = helper.SetError(http.StatusInternalServerError, err.Error())
			return nil, err
		}

		err = helper.SetFilePermissions(srv, files.Id)

		if err != nil{
			fmt.Println("Error setting file permissions:", err.Error())
			err = helper.SetError(http.StatusInternalServerError, err.Error())
			return nil, err
		}

		url := fmt.Sprintf("https://drive.google.com/uc?id=%s", files.Id)

		product.ProductImage = &url
	}

	product.ProductName = productRequest.ProductName
	product.Description = productRequest.Description
	product.Price = productRequest.Price

	product, err = p.productRepository.Update(product)

	if err != nil{
		if productRequest.ProductImage != nil{
			errDelete := helper.DeleteFile(srv, files.Id)
			if errDelete != nil {
				fmt.Println("Error deleting file:", errDelete)
				return nil, errDelete
			}
		}
		return nil, err
	}

	if fileID != ""{
		err = helper.DeleteFile(srv, fileID)
	
		if err != nil {
			fmt.Println("Error deleting file:", err.Error())
			err = helper.SetError(http.StatusInternalServerError, err.Error())
			return nil, err
		}
	}
	
	return product, nil
}

func (p *ProductServiceImpl) Delete(id int) error{

	product,err := p.productRepository.FindByID(id)

	if err != nil{
		return err
	}

	srv, err := helper.GetDriveService()

	if err != nil{
		fmt.Println(err.Error())
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return err
	}

	fileID, err := helper.GetFileIDFromURL(*product.ProductImage)

	if err != nil{
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return err
	}

	err = helper.DeleteFile(srv, fileID)

	if err != nil {
		fmt.Println("Error deleting file:", err.Error())
		err = helper.SetError(http.StatusInternalServerError, err.Error())
		return err
	}

	err = p.productRepository.Delete(id)

	if err != nil{
		return err
	}

	return nil
}