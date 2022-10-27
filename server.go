package main

import (
	"github.com/gofiber/fiber/v2" // gofiber import
)

func main() {
// init fiber app
	app := fiber.New()
// define route
	app.Get("/", func(c *fiber.Ctx) error {
		// return index.html
		return c.SendFile("./public/index.html")
	})
// serve static files
	app.Static("/", "./public")
// start server
	app.Listen(":3000")
}
