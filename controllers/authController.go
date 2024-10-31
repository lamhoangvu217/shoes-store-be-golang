package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
	"github.com/lamhoangvu217/shoes-store-be-golang/utils"
	"log"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func Register(c *fiber.Ctx) error {
	var data map[string]interface{}
	var userData models.User
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("unalble to parse body")
	}
	if len(data["password"].(string)) <= 6 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Password must be greater than 6 character",
		})
	}
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid Email Address",
		})
	}
	// check if email already exist in database
	database.DB.Where("email=?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.ID != 0 {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Email already exist",
		})
	}
	user := models.User{
		Email:    strings.TrimSpace(data["email"].(string)),
		FullName: data["fullName"].(string),
	}
	user.SetPassword(data["password"].(string))
	err := database.DB.Create(&user).Error
	if err != nil {
		log.Println(err)
		c.Status(500)
		return c.JSON(fiber.Map{
			"message": "Could not create account, please try again later.",
		})
	}
	c.Status(200)
	return c.JSON(fiber.Map{
		"user":    user,
		"message": "Account created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		fmt.Println("unalble to parse body")
	}
	var user models.User
	database.DB.Where("email=?", data["email"]).First(&user)
	if user.ID == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Email Address doesn't exist, please create an account",
		})
	}
	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Invalid Password",
		})
	}
	token, err := utils.GenerateJwt(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not log in, please try again later",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message":      "Login successfully",
		"user":         user,
		"access_token": cookie.Value,
	})
}
func Logout(c *fiber.Ctx) error {
	// Create a cookie with the same name as the JWT cookie but set its expiry time in the past
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set expiry to the past to delete the cookie
		HTTPOnly: true,
	}

	// Set the cookie to the response
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}
