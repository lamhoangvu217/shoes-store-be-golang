package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
	"github.com/lamhoangvu217/shoes-store-be-golang/utils"
	"log"
)

func GetUserDetail(c *fiber.Ctx) error {
	// Retrieve JWT token from the cookie
	cookie := c.Cookies("access_token")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse the JWT to get the user email
	userEmail, err := utils.ParseJwt(cookie)
	if err != nil {
		log.Println("Invalid JWT token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token, please log in again",
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
