package helper

import (
	models "belajar-go-restful-api/model/domain"
	"math"
)

func PaginationData(page int64, limit float64, totalRow float64) *models.Metadata{
	totalPage := math.Ceil(float64(totalRow) / limit )
	return &models.Metadata{
		Page: page,
		PerPage: int64(limit),
		TotalRow: int64(totalRow),
		TotalPage: int64(totalPage),
	}
}