package controllers

import (
	"errors"
	"gorm.io/gorm"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
)

var validate = validator.New()

func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DB.Preload("Products").Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":    "get categories successfully",
		"categories": categories,
	})
}

func GetCategoryById(c *fiber.Ctx) error {
	categoryIdStr := c.Params("id")
	if categoryIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "categoryId is required",
		})
	}
	// Convert categoryId from string to uint
	categoryId, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid categoryId",
		})
	}
	var category models.Category
	if err := database.DB.Preload("Products").First(&category, categoryId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":  "get category successfully",
		"category": category,
	})
}

func CreateCategory(c *fiber.Ctx) error {
	category := new(models.Category)
	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}
	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(category)
}

func DeleteCategory(c *fiber.Ctx) error {
	categoryIdStr := c.Params("id")
	if categoryIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "category id is required",
		})
	}
	// Convert product id from string to uint
	categoryId, err := strconv.ParseUint(categoryIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid category id",
		})
	}
	var category models.Category
	if err := database.DB.First(&category, categoryId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "category id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product",
		})
	}
	var productCount int64
	if err := database.DB.Model(&models.Product{}).Where("category_id = ?", categoryId).Count(&productCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if productCount > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category cannot be deleted because it has products now",
		})
	}
	if err := database.DB.Delete(&category, categoryId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "delete category successfully",
	})
}
