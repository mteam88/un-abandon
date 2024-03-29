package app

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/mteam88/un-abandon/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

func GetOauthToken(code string) (token string, err error) {
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")

	// or you can create new url.Values struct and encode that like so
	q := url.Values{}
	q.Add("client_id", os.Getenv("OAUTH_CLIENT_ID"))
	q.Add("client_secret", os.Getenv("OAUTH_CLIENT_SECRET"))
	q.Add("code", code)

	req.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Print(err)
		return "", err
	}

	body := make(map[string]interface{})
	rawBody, err := io.ReadAll(res.Body)
	json.Unmarshal(rawBody, &body)

	if err != nil {
		log.Print(err)
		return "", err
	}
	// return access token
	return body["access_token"].(string), nil
}

func InstallSetup() {
	InstallGroup := App.Group("/install")
	InstallGroup.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("https://github.com/login/oauth/authorize?client_id=064a76c57f88a8d1b666&scope=user,public_repo")
	})
	InstallGroup.Get("/callback", func(c *fiber.Ctx) error {
		// return index.html
		code := c.Query("code")
		token, err := GetOauthToken(code)
		if err != nil {
			log.Print(err)
			return err
		}

		ctx := context.Background()
		client := GetGithubClient(ctx, token)

		// list user info for authenticated user
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			log.Print(err)
			return err
		}

		DB.Update(func(txn *badger.Txn) error {
			// add token to db
			var Users []database.User
			rawusers, err := txn.Get([]byte("users"))
			if err != nil {
				log.Print(err)
				return err
			}
			rawusers.Value(func(val []byte) error {
				json.Unmarshal(val, &Users)
				Users = append(Users, database.User{Username: *user.Login, Token: token, GithubID: *user.ID})
				users, err := json.Marshal(Users)
				if err != nil {
					log.Print(err)
					return err
				}
				txn.Set([]byte("users"), users)
				return nil
			})
			return nil
		})

		// set cookie
		c.Cookie(&fiber.Cookie{
			Name:     "github_token",
			Value:    token,
			HTTPOnly: true,
		})

		// redirect user to home page
		return c.Redirect("/dashboard", 302)
	})
}
