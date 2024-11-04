package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
)

func GetUserDetail(c *fiber.Ctx) error {
	userEmail := c.Locals("userEmail")
	if userEmail == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	// Retrieve user details from the database
	var user models.User
	if err := database.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
			"email":   userEmail,
		})
	}
	return c.JSON(fiber.Map{
		"message": "get user detail successfully",
		"user":    user,
	})
}
