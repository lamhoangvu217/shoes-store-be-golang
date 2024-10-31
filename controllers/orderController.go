package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/helpers"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
	"gorm.io/gorm"
	"net/http"
)

func CreateOrder(c *fiber.Ctx) error {
	order := new(models.Order)
	if err := c.BodyParser(order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order request body",
		})
	}
	// validate product id
	var productId = order.ProductID
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
	// Validate the order status
	if !helpers.IsValidOrderStatus(order.OrderStatus) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Order status is invalid",
		})
	}

	// Validate user
	var userId = order.UserID
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product",
		})
	}

	// validate quantity
	var quantity = order.Quantity
	if quantity < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Quantity is not smaller than 1",
		})
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Order created successfully",
		"order":     order,
		"productId": productId,
	})
}
