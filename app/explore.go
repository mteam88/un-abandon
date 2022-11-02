package app

import (
	"github.com/gofiber/fiber/v2"
)

func ExploreSetup() {
	// define route
	ExploreGroup := App.Group("/explore")

	ExploreGroup.Get("/", func(c *fiber.Ctx) error {
		return c.Render("explore", fiber.Map{
			"Header": "Explore",
		}, "layouts/main")
	})
}