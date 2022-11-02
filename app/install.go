package app

import (
	"net/http"
	"log"
	"net/url"
	"encoding/json"
	"context"
	"os"
	"io"

	"github.com/mteam88/un-abandon/database"

	_ "github.com/joho/godotenv/autoload"
	"github.com/gofiber/fiber/v2"
)

func GetOauthToken(code string) (token string, err error) {
	req, err := http.NewRequest("POST","https://github.com/login/oauth/access_token", nil)
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

	log.Print(req)

	if err != nil {
		log.Print(err)
		return "", err
	}

	body := make(map[string]interface{})
	rawBody,err := io.ReadAll(res.Body)
	json.Unmarshal(rawBody, &body)

	log.Print(body)
	
	if err != nil {
		log.Print(err)
		return "", err
	}
	// return access token
	return body["access_token"].(string), nil
}

func InstallSetup() {
	App.Get("/install", func(c *fiber.Ctx) error {
		// return index.html
		return c.Redirect("https://github.com/login/oauth/authorize?client_id=064a76c57f88a8d1b666&scope=user,public_repo")
	})
	App.Get("/install/callback", func(c *fiber.Ctx) error {
		// return index.html
		code := c.Query("code")
		token, err := GetOauthToken(code)
		if err != nil {
			log.Print(err)
			return err
		}

		log.Print(token)

		ctx := context.Background()
		client := GetGithubClient(ctx, token)

		// list user info for authenticated user
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			log.Print(err)
			return err
		}
		log.Print("user:", *user)

		// add token to db
		var Users []database.User
		rawusers, err := DB.Get("users")
		if err != nil {
			log.Print(err)
			return err
		}
		json.Unmarshal(rawusers, &Users)
		Users = append(Users, database.User{Username: *user.Login, Token: token, GithubID: *user.ID})
		users, err := json.Marshal(Users)
		if err != nil {
			log.Print(err)
			return err
		}
		DB.Set("users", users)

		// set cookie
		c.Cookie(&fiber.Cookie{
			Name: "github_token",
			Value: token,
			HTTPOnly: true,
		})

		// redirect user to home page
		return c.Redirect("/", 302)
	})
}
