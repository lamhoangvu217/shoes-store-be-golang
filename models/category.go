package models

type Category struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"size:255;not null;unique" json:"name"`
	Products []Product `json:"products"`
}
