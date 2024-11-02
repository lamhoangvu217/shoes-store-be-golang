package controllers

import (
	"encoding/json"
	"errors"
	"github.com/lamhoangvu217/shoes-store-be-golang/services"
	"gorm.io/gorm"
	"net/http"
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
	// Call the service to get products by category
	products, err := services.GetProductsByCategory(uint(categoryId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":  "get products successfully",
		"products": products,
	})
}

func GetProductById(c *fiber.Ctx) error {
	productIdStr := c.Params("id")
	if productIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product id is required",
		})
	}
	// Convert categoryId from string to uint
	productId, err := strconv.ParseUint(productIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid categoryId",
		})
	}
	product, err := services.GetProductById(uint(productId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "get product successfully",
		"product": product,
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
			"error": err.Error(),
			"json":  product,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	productIdStr := c.Params("id")
	if productIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product id is required",
		})
	}
	// Convert product id from string to uint
	productId, err := strconv.ParseUint(productIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid productId",
		})
	}
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product",
		})
	}

	if err := database.DB.Delete(&product, productId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "delete product successfully",
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	productIdStr := c.Params("id")
	if productIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product id is required",
		})
	}
	// Convert product id from string to uint
	productId, err := strconv.ParseUint(productIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid productId",
		})
	}
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product",
		})
	}
	var updateProductData models.Product
	if err := c.BodyParser(&updateProductData); err != nil {
		// Return 400 if request body is invalid
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if updateProductData.Name == "" || len(updateProductData.Name) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product name is required and must be at least 3 characters long",
		})
	}
	if updateProductData.OriginalPrice <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Original price must be greater than zero",
		})
	}
	if updateProductData.SalePrice <= 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Sale price must be greater than zero",
		})
	}
	if updateProductData.CategoryID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Category ID is required",
		})
	}
	// Update the product fields
	product.Name = updateProductData.Name
	product.OriginalPrice = updateProductData.OriginalPrice
	product.SalePrice = updateProductData.SalePrice
	product.Description = updateProductData.Description
	product.CategoryID = updateProductData.CategoryID
	product.ImageUrl = updateProductData.ImageUrl
	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

func CreateProductImagesGallery(c *fiber.Ctx) error {
	productImage := new(models.ProductImage)
	if err := c.BodyParser(productImage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}
	if err := database.DB.Create(&productImage).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(productImage)
}

func DeleteProductImage(c *fiber.Ctx) error {
	productImgIdStr := c.Params("id")
	if productImgIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product id is required",
		})
	}
	// Convert product id from string to uint
	productImgId, err := strconv.ParseUint(productImgIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product image id",
		})
	}
	var productImg models.ProductImage
	if err := database.DB.First(&productImg, productImgId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product image id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product",
		})
	}

	if err := database.DB.Delete(&productImg, productImgId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "delete product image successfully",
	})
}
func UpdateProductImage(c *fiber.Ctx) error {
	productImgIdStr := c.Params("id")
	if productImgIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product image id is required",
		})
	}
	// Convert product id from string to uint
	productImgId, err := strconv.ParseUint(productImgIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product image id",
		})
	}
	var productImg models.ProductImage
	if err := database.DB.First(&productImg, productImgId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product image id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product image",
		})
	}
	var updateProductImageData models.ProductImage
	if err := c.BodyParser(&updateProductImageData); err != nil {
		// Return 400 if request body is invalid
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if updateProductImageData.Url == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "url is required",
		})
	}
	productImg.Url = updateProductImageData.Url
	productImg.ProductID = updateProductImageData.ProductID
	// check product id exist in products table
	if err := database.DB.Save(&productImg).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "product image updated successfully",
		"productImage": productImg,
	})
}

func GetProductImages(c *fiber.Ctx) error {
	var productImages []models.ProductImage
	if err := database.DB.Find(&productImages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":       "get product images successfully",
		"productImages": productImages,
	})
}
func GetProductImagesByProductId(c *fiber.Ctx) error {
	productIdStr := c.Query("productId")
	if productIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "productId is required",
		})
	}
	// Convert categoryId from string to uint
	productId, err := strconv.ParseUint(productIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid productId",
		})
	}
	// check product id exist in database
	var product models.Product
	if err := database.DB.First(&product, productId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product",
		})
	}
	var productImages []models.ProductImage
	if err := database.DB.Where("product_id = ?", productId).Preload("Product").Find(&productImages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":       "get product images successfully",
		"productImages": productImages,
	})
}
