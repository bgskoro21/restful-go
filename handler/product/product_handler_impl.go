package handler

import (
	"belajar-go-restful-api/helper"
	"belajar-go-restful-api/model/web"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gopkg.in/h2non/filetype.v1"
)

func isValidString(input string) bool{
	trimmed := strings.TrimSpace(input)
	return len(trimmed) > 0
}

func isValidImage(file *multipart.FileHeader) error{
	buffer := make([]byte, 512)
	maxSize := int64(2 << 20)

	if file.Size > maxSize{
		err := helper.SetError(http.StatusBadRequest, "File size too large!")
		return err
	}

	src,err := file.Open()

	if err != nil{
		err = helper.SetError(http.StatusBadRequest, "File cannot be accessed!")
		return err
	}

	defer src.Close()

	_, err = src.Read(buffer)

	if err != nil{
		err = helper.SetError(http.StatusBadRequest, "File cannot be read!")
		return err
	}

	kind,_ := filetype.Match(buffer)

	if kind.MIME.Value == "image/jpeg" || kind.MIME.Value == "image/png" {
		return nil
	}

	err = helper.SetError(http.StatusBadRequest, "File must be png or jpg!")
	return err
}

func (h *ProductHandlerImpl) FindAll(c *gin.Context){
	var request web.ProductGetRequest

	if err := c.ShouldBindQuery(&request); err != nil{
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors){
			errorMessage := fmt.Sprintf("Error on filed %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	products, metadata, _ := h.productService.FindAll(request)

	c.JSON(http.StatusOK, gin.H{
		"Code": 200,
		"Status": "OK",
		"data": products,
		"pagination": metadata,
	})
}

func (h *ProductHandlerImpl) FindByID(c *gin.Context){
	id,_ := strconv.ParseInt(c.Param("id"), 10, 64)

	if id == 0{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID Invalid",
		})

		return
	}

	product, err := h.productService.FindByID(int(id))

	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{
			"Code": 200,
			"Status": "Not Found",
			"Message": "Data product not found!",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Code": 200,
		"Status": "OK",
		"data": product,
	})
}

func (h *ProductHandlerImpl) Create(c *gin.Context){
	var request web.ProductCreateRequest

	if err := c.ShouldBind(&request); err != nil{
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			errorMessage := fmt.Sprintf("Error in field %s, condition must be %s", e.Field, e.Type)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMessage,
			})
		case validator.ValidationErrors:
			errorMessages := []string{}
			for _, err := range e {
				errorMessage := fmt.Sprintf("Error on field %s, condition: %s", err.Field(), err.ActualTag())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMessages,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
		return
	}

	if !isValidString(request.ProductName){
		c.JSON(http.StatusBadRequest, gin.H{
			"Code": 400,
			"Status": "Bad Request",
			"Message": "ProductName is not valid!",
		})
		return
	}

	if err := isValidImage(request.ProductImage); err != nil{
		errResponse := helper.GetError(err)
		c.JSON(int(errResponse.Code()), gin.H{
			"Code": errResponse.Code(),
			"Status": "Bad Request",
			"Message": err.Error(),
		})
		return
	}

	product, err := h.productService.Create(request)

	if err != nil{
		errorResponse := helper.GetError(err)

		c.JSON(int(errorResponse.Code()), gin.H{
			"error": errorResponse.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Code": 200,
		"Status": "OK",
		"Data": product,
	})
}

func (h *ProductHandlerImpl) Update(c *gin.Context){
	var request web.ProductUpdateRequest

	if err := c.ShouldBind(&request); err != nil{
		switch e := err.(type) {
		case *json.UnmarshalTypeError:
			errorMessage := fmt.Sprintf("Error in field %s, condition must be %s", e.Field, e.Type)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMessage,
			})
		case validator.ValidationErrors:
			errorMessages := []string{}
			for _, err := range e {
				errorMessage := fmt.Sprintf("Error on field %s, condition: %s", err.Field(), err.ActualTag())
				errorMessages = append(errorMessages, errorMessage)
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errorMessages,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	
	request.ID = id

	product, err := h.productService.Update(request)

	if err != nil{
		errorResponse := helper.GetError(err)
		
		c.JSON(int(errorResponse.Code()), gin.H{
			"error": errorResponse.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Code": 200,
		"Status": "OK",
		"Data": product,
	})
}

func (h *ProductHandlerImpl) Delete(c *gin.Context){
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.productService.Delete(id)

	if err != nil{
		errorResponse := helper.GetError(err)
		c.JSON(int(errorResponse.Code()), gin.H{
			"error": errorResponse.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Code": 200,
		"Status": "OK",
		"Message": "Data product successful deleted!",
	})
}

func (h *ProductHandlerImpl) GetProductImage(c *gin.Context){
	filename := c.Param("filename")

	filepath := "storage/" + filename

	c.File(filepath)
}