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
	admin.Put("/product/:id", controllers.UpdateProduct)
	admin.Put("/category/:id", controllers.UpdateCategory)

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/logout", controllers.Logout)

	app.Post("/api/create-order", controllers.CreateOrder)
}
