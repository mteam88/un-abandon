package app

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Repo struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Url         string `json:"html_url"`
}

func DashboardSetup() {
	App.Get("/dashboard", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		repos, _, err := client.Repositories.List(ctx, "", nil)

		if err != nil {
			log.Print(err)
			return err
		}
		// clean repos object to only include name, url and description
		var cleanRepos []Repo = []Repo{}
		for _, repo := range repos {
			cleanRepos = append(cleanRepos, Repo{
				Name:        repo.GetName(),
				Description: repo.GetDescription(),
				Url:         repo.GetHTMLURL(),
			})
		}

		return c.Render("dashboard", fiber.Map{
			"Header": "Dashboard",
			"Repos": cleanRepos,
		}, "layouts/main")
	})
}