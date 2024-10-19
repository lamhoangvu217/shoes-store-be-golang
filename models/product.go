package models

type Product struct {
	ID            uint     `gorm:"primaryKey" json:"id"`
	Name          string   `gorm:"size:255;not null" json:"name"`
	OriginalPrice float64  `gorm:"not null" json:"originalPrice"`
	SalePrice     float64  `gorm:"not null" json:"salePrice"`
	Description   string   `gorm:"type:text" json:"description"`
	CategoryID    uint     `gorm:"not null" json:"categoryId"`
	Category      Category `json:"category"`
}
