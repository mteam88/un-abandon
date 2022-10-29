package main

import (
	"github.com/gofiber/fiber/v2" // gofiber import
	"github.com/gofiber/template/html" // mustache template import
)

func main() {
// init fiber app
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"), // set the views directory
	})
// define route
	app.Get("/", func(c *fiber.Ctx) error {
		// return index.html
		return c.Render("index", fiber.Map{
			"Header": "Un-Abandon",
		}, "layouts/main")
	})
	app.Get("/explore", func(c *fiber.Ctx) error {
		// return index.html
		return c.Render("explore", fiber.Map{
			"Header": "Explore",
		}, "layouts/main")
	})
	app.Get("/install", func(c *fiber.Ctx) error {
		// return index.html
		return c.Render("install", fiber.Map{
			"Header": "Install",
		}, "layouts/main")
	})
// serve static files
	app.Static("/", "./static/public")
// start server
	app.Listen(":3000")
}
