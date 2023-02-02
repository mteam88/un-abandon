package app

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/mteam88/un-abandon/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/go-github/v48/github"
)

func DashboardSetup() {
	DashboardGroup := App.Group("/dashboard", AuthenticateUser)

	DashboardGroup.Get("/", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		repos, _, err := client.Repositories.List(ctx, "", nil)

		if err != nil {
			log.Print(err)
			return err
		}

		abandonedRepos, err := GetAllAbandonedRepos()
		if err != nil {
			log.Print(err)
			return err
		}

		// clean repos object to only include id, name, url and description
		var cleanRepos []database.Repo = []database.Repo{}
		for _, repo := range repos {
			// check if repo.ID is in abandonedRepos
			var found bool = false
			for _, abandonedRepo := range abandonedRepos {
				if repo.GetID() == abandonedRepo.ID {
					found = true
					break
				}
			}
			if !found {
				cleanRepos = append(cleanRepos, database.Repo{
					Name:        repo.GetName(),
					Description: repo.GetDescription(),
					Url:         repo.GetHTMLURL(),
					ID:          repo.GetID(),
				})
			}
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

	DashboardGroup.Post("/abandon/", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		var id struct {
			ID int64 `json:"id"`
		}

		err := json.Unmarshal(c.Body(), &id)
		if err != nil {
			log.Print(err)
			return err
		}

		repo, _, err := client.Repositories.GetByID(ctx, id.ID)
		if err != nil {
			log.Print(err)
			return err
		}

		// add repo to db
		err = RepoDB.Update(func(txn *badger.Txn) error {
			err := txn.Set(database.GetID(int32(repo.GetID())), database.Repo{
				Name:        repo.GetName(),
				Description: repo.GetDescription(),
				Url:         repo.GetHTMLURL(),
				ID:          repo.GetID(),
				Token:       c.Cookies("github_token"),
			}.Serialize())
			return err
		})
		if err != nil {
			log.Print(err)
			return err
		}
		return c.SendStatus(200)
	})

	DashboardGroup.Post("/adopt/", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		var id struct {
			ID int32 `json:"id"`
		}

		err := json.Unmarshal(c.Body(), &id)
		if err != nil {
			log.Print(err)
			return err
		}

		err = RepoDB.Update(func(txn *badger.Txn) error {
			item, err := txn.Get(database.GetID(id.ID))
			if err != nil {
				return err
			}
			item.Value(func(val []byte) error {
				dbrepo := database.DeserializeRepo(val)
				err = txn.Delete(database.GetID(id.ID))
				if err != nil {
					return err
				}
				// transfer repo
				err = TransferRepo(dbrepo, client)
				if err != nil {
					return err
				}
				return nil
			})
			return err
		})
		if err != nil {
			log.Print(err)
			return err
		}
		// return ok
		return c.SendStatus(200)
	})

	// log out route
	App.Get("/logout/", func(c *fiber.Ctx) error {
		c.ClearCookie()
		return c.Redirect("/")
	})
}

func TransferRepo(dbrepo database.Repo, newOwnerClient *github.Client) error {
	ctx := context.Background()
	currentOwnerClient := GetGithubClient(ctx, dbrepo.Token)

	// get repo by id
	ghrepo, _, err := currentOwnerClient.Repositories.GetByID(ctx, dbrepo.ID)
	if err != nil {
		log.Print(err)
		return err
	}

	// get new client username
	newOwner, _, err := newOwnerClient.Users.Get(ctx, "")
	if err != nil {
		log.Print(err)
		return err
	}
	newOwnerUsername := newOwner.GetLogin()
	// transfer repo
	_, _, err = currentOwnerClient.Repositories.Transfer(ctx, ghrepo.GetOwner().GetLogin(), ghrepo.GetName(), github.TransferRequest{
		NewOwner: newOwnerUsername,
	})
	if err != nil {
		if strings.Contains(err.Error(), "Repositories cannot be transferred to the original owner") {
			// repo already transferred
			return errors.New("repo cannot be transferred to original owner")
		} else if strings.Contains(err.Error(), "job scheduled on GitHub side; try again later") {
			log.Print(err)
			return nil
		} else if strings.Contains(err.Error(), "Repository has already been taken") {
			log.Print(err)
			return err
		} else {
			log.Print(err)
			return err
		}
	}

	return nil
}

func GetAllAbandonedRepos() ([]database.Repo, error) {
	var repos []database.Repo
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
			var repo = database.DeserializeRepo(v)
			repos = append(repos, repo)
		}
		return nil
	})
	return repos, err
}
