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
		
		err := RepoDB.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = true
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				v, err := item.ValueCopy(nil)
				if err != nil {
					return err
				}
				var repo database.Repo
				err = json.Unmarshal(v, &repo)
				if err != nil {
					return err
				}
				// remove sensitive fields from repo
				repo.Token = ""
				cleanRepos = append(cleanRepos, repo)
			}
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
