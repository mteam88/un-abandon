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

		// clean repos object to only include name, url and description
		var cleanRepos []database.Repo = []database.Repo{}
		for _, repo := range repos {
			// check if repo.ID is in abandonedRepos
			var found bool = false
			for _, abandonedRepo := range abandonedRepos {
				if repo.GetID() == abandonedRepo {
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

		var url struct {
			Url string `json:"url"`
		}

		err := json.Unmarshal(c.Body(), &url)
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
				err := DB.Update(func(txn *badger.Txn) error {
					abandoned_repos, err := txn.Get([]byte("abandoned_repos"))
					if err != nil {
						log.Print(err)
					}
					abandoned_repos.Value(func(val []byte) error {
						var repos []database.Repo
						err = json.Unmarshal(val, &repos)
						if err != nil {
							log.Print(err)
							return err
						}
						repos = append(repos, database.Repo{
							Name:        repo.GetName(),
							Description: repo.GetDescription(),
							Url:         repo.GetHTMLURL(),
							ID:          repo.GetID(),
							Token:       c.Cookies("github_token"),
						})
						new_abandoned_repos, err := json.Marshal(repos)
						if err != nil {
							log.Print(err)
							return err
						}
						err = txn.Set([]byte("abandoned_repos"), new_abandoned_repos)
						if err != nil {
							log.Print(err)
							return err
						}
						return err
					})
					return err
				})
				if err != nil {
					log.Print(err)
					return err
				}
			}
		}
		return c.Redirect("/dashboard")
	})

	DashboardGroup.Post("/adopt/", func(c *fiber.Ctx) error {
		ctx := context.Background()
		client := GetGithubClient(ctx, c.Cookies("github_token"))

		var url struct {
			Url string `json:"url"`
		}

		err := json.Unmarshal(c.Body(), &url)
		if err != nil {
			log.Print(err)
			return err
		}

		err = DB.Update(func(txn *badger.Txn) error {
			abandoned_repos, err := txn.Get([]byte("abandoned_repos"))
			if err != nil {
				log.Print(err)
				return err
			}
			err = abandoned_repos.Value(func(val []byte) error {
				var repos []database.Repo
				err = json.Unmarshal(val, &repos)
				if err != nil {
					log.Print(err)
					return err
				}

				for i, repo := range repos {
					if repo.Url == url.Url {
						err = TransferRepo(repo, client)
						if err != nil {
							log.Print(err)
							return err
						}
						// remove repo from db
						repos = append(repos[:i], repos[i+1:]...)
						new_abandoned_repos, err := json.Marshal(repos)
						if err != nil {
							log.Print(err)
							return err
						}
						err = txn.Set([]byte("abandoned_repos"), new_abandoned_repos)
						if err != nil {
							log.Print(err)
							return err
						}
					}
				}
				return err
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

func GetAllAbandonedRepos() ([]int64, error) {
	var id_abandoned_repos []int64
	err := DB.View(func(txn *badger.Txn) error {
	abandoned_repos, err := txn.Get([]byte("abandoned_repos"))
	if err != nil {
		log.Print(err)
		return err
	}
	abandoned_repos.Value(func(val []byte) error {
	var repos []database.Repo
	err = json.Unmarshal(val, &repos)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, repo := range repos {
		id_abandoned_repos = append(id_abandoned_repos, repo.ID)
	}
	return nil
})
return nil
})
return id_abandoned_repos, err
}
