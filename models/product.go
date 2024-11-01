package models

type Product struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:255;not null" json:"name" validate:"required,min=3"`
	OriginalPrice float64        `gorm:"not null" json:"originalPrice" validate:"required,gt=0"`
	SalePrice     float64        `gorm:"not null" json:"salePrice" validate:"required,gt=0"`
	Description   string         `gorm:"type:text" json:"description"`
	ImageUrl      string         `gorm:"type:text" json:"imageUrl"`
	ProductImages []ProductImage `json:"productImages"`
	CategoryID    uint           `gorm:"not null" json:"categoryId"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"-"`
}
type ProductImage struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Url       string  `gorm:"type:text;not null" json:"url"`
	ProductID uint    `gorm:"not null" json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"-"`
}
