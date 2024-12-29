package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id"`
	User       User        `json:"-"`
	Status     string      `json:"status"`
	Total      float64     `json:"total"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"-"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
