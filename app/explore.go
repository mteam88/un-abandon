package app

import (
	"github.com/gofiber/fiber/v2"
)

func ExploreSetup() {
	// define route
	ExploreGroup := App.Group("/explore")

	ExploreGroup.Get("/", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		repos, _, err := client.Repositories.List(ctx, "", nil)

		if err != nil {
			log.Print(err)
			return err
		}
		// clean repos object to only include name, url and description
		var cleanRepos []database.Repo = []database.Repo{}
		for _, repo := range repos {
			cleanRepos = append(cleanRepos, database.Repo{
				Name:        repo.GetName(),
				Description: repo.GetDescription(),
				Url:         repo.GetHTMLURL(),
				ID:          repo.GetID(),
			})
		}

		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			log.Print(err)
			return err
		}

		return c.Render("explore", fiber.Map{
			"Header": "Explore",
		}, "layouts/main")
	})
}