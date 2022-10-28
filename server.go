package main

import (
	"github.com/gofiber/fiber/v2" // gofiber import
	"github.com/gofiber/template/mustache" // mustache template import
)

func main() {
// init fiber app
	app := fiber.New(fiber.Config{
		Views: mustache.New("./views", ".mustache"),
	})
// define route
	app.Get("/", func(c *fiber.Ctx) error {
		// return index.html
		return c.Render("index", fiber.Map{})
	})
	app.Get("/explore", func(c *fiber.Ctx) error {
		// return index.html
		return c.Render("explore", fiber.Map{})
	})
// serve static files
	app.Static("/", "./static/public")
// start server
	app.Listen(":3000")
}
