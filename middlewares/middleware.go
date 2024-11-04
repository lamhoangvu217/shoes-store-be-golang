package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/utils"
	"log"
)

func AuthRequired(c *fiber.Ctx) error {
	// Get the access token cookie
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

	// Store the userEmail in the request context for future use
	c.Locals("userEmail", userEmail)

	// Continue to the next handler
	return c.Next()
}
