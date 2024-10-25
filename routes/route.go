package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lamhoangvu217/shoes-store-be-golang/controllers"
)

func Setup(app *fiber.App) {
	app.Get("/api/categories", controllers.GetCategories)
	app.Get("/api/products", controllers.GetProductsByCategory)
	app.Get("/api/category/:id", controllers.GetCategoryById)

	admin := app.Group("/admin")
	admin.Post("/product", controllers.CreateProduct)
	admin.Post("/category", controllers.CreateCategory)
	admin.Delete("/product/:id", controllers.DeleteProduct)
	admin.Delete("/category/:id", controllers.DeleteCategory)

}
