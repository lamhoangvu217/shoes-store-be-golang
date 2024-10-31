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
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Parse the JWT to get the user ID
	userId, err := utils.ParseJwt(cookie)
	if err != nil {
		log.Println("Invalid JWT token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token, please log in again",
		})
	}
	// Convert userId to integer for database lookup
	//userIdInt, err := strconv.Atoi(userId)
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"message": "Invalid user ID",
	//		"userId":  userIdInt,
	//	})
	//}
	// Retrieve user details from the database
	var user models.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	return c.JSON(fiber.Map{
		"message": "get user detail successfully",
		"user":    user,
	})
}
