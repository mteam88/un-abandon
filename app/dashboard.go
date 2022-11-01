package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mteam88/un-abandon/database"

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

		return c.Render("dashboard", fiber.Map{
			"Header":   "Dashboard",
			"Repos":    cleanRepos,
			"Username": user.Login,
		}, "layouts/main")
	})

	App.Post("/dashboard/abandon/", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		var url struct {
			Url string `json:"url"`
		}

		err := json.Unmarshal(c.Body(), url)
		if err != nil {
			log.Print(err)
			return err
		}

		repos, _, err := client.Repositories.List(ctx, "", nil)
		if err != nil {
			log.Print(err)
			return err
		}

		for _, repo := range repos {
			if repo.GetHTMLURL() == url.Url {
				// add repo to db
				abandoned_repos, err := DB.Get("abandoned_repos")
				if err != nil {
					log.Print(err)
					return err
				}
				var repos []database.Repo
				err = json.Unmarshal(abandoned_repos, &repos)
				if err != nil {
					log.Print(err)
					return err
				}
				repos = append(repos, database.Repo{
					Name:        repo.GetName(),
					Description: repo.GetDescription(),
					Url:         repo.GetHTMLURL(),
					ID:          repo.GetID(),
				})
				abandoned_repos, err = json.Marshal(repos)
				if err != nil {
					log.Print(err)
					return err
				}
				err = DB.Set("abandoned_repos", abandoned_repos)
				if err != nil {
					log.Print(err)
					return err
				}
			}
		}

		return c.Redirect("/dashboard")
	})

}
