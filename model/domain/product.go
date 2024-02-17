package models

import "time"

type Product struct{
	ID uint `gorm:"primaryKey" json:"id"`
	ProductName string `gorm:"type:varchar(100)" json:"product_name"`
	Description string `json:"description"`
	Price int `json:"price"`
	ProductImage *string `json:"product_image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}