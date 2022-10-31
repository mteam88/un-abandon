package app

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

func DashboardSetup() {
	App.Get("/dashboard", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		repos, _, err := client.Repositories.List(ctx, "", nil)

		if err != nil {
			log.Print(err)
			return err
		}

		return c.Render("dashboard", fiber.Map{
			"Header": "Dashboard",
			"Repos": repos,
		}, "layouts/main")
	})
}