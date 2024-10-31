package models

import "time"

type Order struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"not null" json:"userId"`                             // Foreign key to the User model
	ProductID   uint      `gorm:"not null" json:"productId"`                          // Foreign key to the Product model
	Quantity    uint      `gorm:"not null" json:"quantity" validate:"required,min=1"` // Number of items ordered
	TotalPrice  float64   `gorm:"not null" json:"totalPrice"`                         // Total price for the order
	OrderStatus string    `gorm:"size:50;not null" json:"orderStatus"`                // Status of the order (e.g., Pending, Completed, Cancelled)
	CreatedAt   time.Time `json:"createdAt"`                                          // Timestamp for when the order was created
	UpdatedAt   time.Time `json:"updatedAt"`                                          // Timestamp for when the order was last updated
}
