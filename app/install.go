package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func InstallSetup() {
	App.Get("/install", func(c *fiber.Ctx) error {
		fmt.Println(c.AllParams()) 
		// return index.html
		return c.Render("install", fiber.Map{
			"Header": "Install",
		}, "layouts/main")
	})
	App.Post("/install", func(c *fiber.Ctx) error {
		// return index.html
		return c.Render("install", fiber.Map{
			"Header": "Install",
		}, "layouts/main")
	})
}
