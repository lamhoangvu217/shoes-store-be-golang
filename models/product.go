package models

type Product struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	Name          string  `gorm:"size:255;not null" json:"name" validate:"required,min=3"`
	OriginalPrice float64 `gorm:"not null" json:"originalPrice" validate:"required,gt=0"`
	SalePrice     float64 `gorm:"not null" json:"salePrice" validate:"required,gt=0"`
	Description   string  `gorm:"type:text" json:"description"`
	CategoryID    uint    `gorm:"not null" json:"categoryId"`
}
