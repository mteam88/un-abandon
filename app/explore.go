package app

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mteam88/un-abandon/database"
)

func ExploreSetup() {
	// define route
	ExploreGroup := App.Group("/explore")

	ExploreGroup.Get("/", func(c *fiber.Ctx) error {
		// get all repos from database
		rawRepos, err := DB.Get("abandoned_repos")
		if err != nil {
			log.Print(err)
			return err
		}

		var repos []database.Repo
		json.Unmarshal(rawRepos, &repos)

		// clean repos object to only include name, url and description
		var cleanRepos []database.Repo = []database.Repo{}
		for _, repo := range repos {
			cleanRepos = append(cleanRepos, database.Repo{
				Name:        repo.Name,
				Description: repo.Description,
				Url:         repo.Url,
				ID:          repo.ID,
			})
		}

		return c.Render("explore", fiber.Map{
			"Header": "Explore",
			"Repos":  cleanRepos,
		}, "layouts/main")
	})
}
