package app

import (
	"log"
	"context"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v48/github"
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2/middleware/session"
)

func AuthSetup() {
	// define middleware
	App.Use(AuthenticateUser)
}

func AuthenticateUser(c *fiber.Ctx) error {
	// check that user is authenticated to github
	// if not, redirect to /install
	log.Print("Authenticating user: " + c.Path() + " token: " + c.Cookies("github_token"))
	if c.Path() == "/dashboard" {
		if c.Cookies("github_token") == "" {
			return c.Redirect("/install")
		} else {
			if (CheckGHOauthToken(c.Cookies("github_token"))) {
				return c.Next()
			}
			return c.Redirect("/install")
		}
	}
	c.Next()
	return nil
}

func CheckGHOauthToken(token string) (ok bool) {
	// check that token is valid
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// check that client is authenticated
	_, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Print("Oauth check failed", err)
		return false
	}
	log.Print("Oauth check passed")
	return true
}