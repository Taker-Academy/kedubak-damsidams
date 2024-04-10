package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	//register
	app.Options("/auth/register", func(c *fiber.Ctx) error {
        return c.SendStatus(fiber.StatusOK)
    })
	app.Post("/auth/register", controllers.CreateUser)

	//login
	app.Options("/auth/login", func(c *fiber.Ctx) error {
        return c.SendStatus(fiber.StatusOK)
    })
	app.Post("/auth/login", controllers.Login)

	//get infos
	app.Options("/user/me", func(c *fiber.Ctx) error {
        return c.SendStatus(fiber.StatusOK)
    })
	app.Get("/user/me", controllers.GetCurrentUser)

}
