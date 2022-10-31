package app

import (
	"context"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v48/github"
)

func GetGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client
}