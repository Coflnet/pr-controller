package github

import (
	"context"
	"time"

	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/rs/zerolog/log"
)

func ListPrs() ([]*model.Pr, error) {

	owner := "Coflnet"
	repo := "hypixel-react"

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	prs, _, err := client.PullRequests.List(ctx, owner, repo, nil)

	if err != nil {
		log.Error().Err(err).Msgf("there was a problem when listing pull requests")
		return nil, err
	}

	log.Info().Msgf("read %d pull requests from github", len(prs))

	var mongoPrs []*model.Pr
	for _, githubPr := range prs {
		pr := model.Pr{
			Owner:      *githubPr.Head.Repo.Owner.Login,
			Repo:       *githubPr.Head.Repo.Name,
			Number:     *githubPr.Number,
			Branch:     *githubPr.Head.Ref,
			LastCommit: *githubPr.Head.SHA,
			Image:      "coflnet/pr-env",
		}

		mongoPrs = append(mongoPrs, &pr)
	}

	return mongoPrs, nil
}
