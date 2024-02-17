package models

type Metadata struct{
	Page int64 `json:"page"`
	TotalRow int64 `json:"total_row"`
	TotalPage int64 `json:"total_page"`
	PerPage int64 `json:"per_page"`
}