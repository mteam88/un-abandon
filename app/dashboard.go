package app

import (
	"github.com/gofiber/fiber/v2"
)

func DashboardSetup() {
	App.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.Render("dashboard", fiber.Map{
			"Header": "Dashboard",
		}, "layouts/main")
	})
}