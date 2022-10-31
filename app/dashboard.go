package app

import (
	"context"
	"log"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v48/github"
	"github.com/gofiber/fiber/v2"
)

func DashboardSetup() {
	App.Get("/dashboard", func(c *fiber.Ctx) error {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.Cookies("github_token")},
		)
		log.Print("token: ", c.Cookies("github_token"))
		tc := oauth2.NewClient(ctx, ts)
	
		client := github.NewClient(tc)

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