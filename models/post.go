package models

type Post struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:255;not null" json:"name" validate:"required,min=3"`
	Description string `gorm:"type:text" json:"description"`
}
