package models

import "time"

type Order struct {
	ID          int         `json:"id"`
	UserID      string        `json:"user_id"`
	Status      OrderStatus      `json:"status"`
	TotalAmount float64     `json:"total_amount"`
	Items       []OrderItem `json:"item,omitempty"`
	CreateAt    time.Time   `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


type OrderItem struct{
	ID int `json:"id"`
	OrderID int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`

	// Optional product info
	ProductName string `json:"product_name,omitempty"`
}


