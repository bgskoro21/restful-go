package helper

import (
	"net/http"

	"github.com/go-openapi/errors"
)


func PanicIfError(err error){
	if err != nil {
		panic(err)
	}
}

func SetError(code int, message string) error{
	return errors.New(int32(code), message)
}

func GetError(err error) errors.Error{
	if v,ok := err.(errors.Error); ok{
		return v
	}

	return errors.New(http.StatusInternalServerError, err.Error())
}