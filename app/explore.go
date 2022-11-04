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
		repos, err := DB.Get("abandoned_repos")
		if err != nil {
			log.Print(err)
			return err
		}

		var cleanRepos []database.Repo
		json.Unmarshal(repos, &cleanRepos)

		log.Print(cleanRepos)

		return c.Render("explore", fiber.Map{
			"Header": "Explore",
			"Repos":  cleanRepos,
		}, "layouts/main")
	})
}
