package services

import (
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
)

func GetProductsByCategory(categoryId uint) ([]models.Product, error) {
	var products []models.Product
	if err := database.DB.Where("category_id = ?", categoryId).Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
func GetProductById(productId uint) (*models.Product, error) {
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		return nil, err
	}
	return &product, nil
}
