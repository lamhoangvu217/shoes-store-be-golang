package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/database"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetPosts get list posts
func GetPosts(c *fiber.Ctx) error {
	var posts []models.Post
	if err := database.DB.Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "get posts successfully",
		"posts":   posts,
	})
}

func GetPostById(c *fiber.Ctx) error {
	postIdStr := c.Params("id")
	if postIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "categoryId is required",
		})
	}
	// Convert categoryId from string to uint
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid categoryId",
		})
	}
	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "get post successfully",
		"post":    post,
	})
}
func CreatePost(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input data",
		})
	}
	if err := database.DB.Create(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(post)
}
func DeletePost(c *fiber.Ctx) error {
	postIdStr := c.Params("id")
	if postIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "post id is required",
		})
	}
	// Convert product id from string to uint
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid post id",
		})
	}
	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "post id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve post",
		})
	}

	if err := database.DB.Delete(&post, postId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "delete post successfully",
	})
}
func UpdatePost(c *fiber.Ctx) error {
	postIdStr := c.Params("id")
	if postIdStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "post id is required",
		})
	}
	// Convert product id from string to uint
	postId, err := strconv.ParseUint(postIdStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid post id",
		})
	}
	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "post id not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve product",
		})
	}
	var updatePostData models.Post
	if err := c.BodyParser(&updatePostData); err != nil {
		// Return 400 if request body is invalid
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	if updatePostData.Name == "" || len(updatePostData.Name) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post name is required and must be at least 3 characters long",
		})
	}
	post.Name = updatePostData.Name
	if err := database.DB.Save(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Post updated successfully",
		"product": post,
	})
}
