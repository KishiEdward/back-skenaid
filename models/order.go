package models

import "time"

type Order struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	UserID          uint        `json:"user_id"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `gorm:"default:'pending'" json:"status"`
	ShippingAddress string      `json:"shipping_address"`
	Notes           string      `json:"notes"`
	PaymentMethod   string      `json:"payment_method"`
	Items           []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	OrderID     uint    `json:"order_id"`
	ProductID   uint    `json:"product_id"`
	Product     Product `gorm:"foreignKey:ProductID" json:"product"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}