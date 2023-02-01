package app

import (
	"encoding/json"
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/mteam88/un-abandon/database"
)

func ExploreSetup() {
	// define route
	ExploreGroup := App.Group("/explore")
	ExploreGroup.Use(AuthenticateUser)

	ExploreGroup.Get("/", func(c *fiber.Ctx) error {
		// get all repos from database
		var cleanRepos []database.Repo = []database.Repo{}
		err := DB.View(func(txn *badger.Txn) error {
		rawRepos, err := txn.Get([]byte("abandoned_repos"))
		if err != nil {
			log.Print(err)
			return err
		}

		rawRepos.Value(func(val []byte) error {
		var repos []database.Repo
		json.Unmarshal(val, &repos)

		// clean repos object to only include name, url and description
		for _, repo := range repos {
			cleanRepos = append(cleanRepos, database.Repo{
				Name:        repo.Name,
				Description: repo.Description,
				Url:         repo.Url,
				ID:          repo.ID,
			})
		}
		return nil
		})
		return nil
		})
		if err != nil {
			log.Print(err)
			return err
		}
		return c.Render("explore", fiber.Map{
			"Header": "Explore",
			"Repos":  cleanRepos,
		}, "layouts/main")
	})
}
