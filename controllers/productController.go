package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/helpers"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
)

var formValidation = validator.New()

func GetProductsByCategory(c *fiber.Ctx) error {
	// Extract categoryId from query parameters
	categoryIdStr := c.Query("categoryId")
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
	var products []models.Product
	// Fetch products belonging to the specified category
	if err := database.DB.Where("category_id = ?", categoryId).Preload("Category").Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":  "get products successfully",
		"products": products,
	})

}

func CreateProduct(c *fiber.Ctx) error {
	product := new(models.Product)

	// Parse the request body into the product struct
	// Attempt to parse the request body into the product struct
	if err := json.Unmarshal(c.Body(), product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON: " + err.Error(),
		})
	}
	err := formValidation.Struct(product)
	if err != nil {
		// If validation fails, return a detailed error message
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = helpers.GetValidationErrorMessage(err)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": errors,
		})
	}
	if err := database.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(product)
}
