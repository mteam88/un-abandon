package app

import (
	"github.com/gofiber/fiber/v2"
)

func ExploreSetup() {
	// define route
	App.Get("/explore", func(c *fiber.Ctx) error {
		return c.Render("explore", fiber.Map{
			"Header": "Explore",
		}, "layouts/main")
	})
}