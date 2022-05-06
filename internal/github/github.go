package github

import (
	"context"
	"os"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"
)

var client *github.Client

func Init() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
}
