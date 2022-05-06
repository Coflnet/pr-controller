package usecase

import (
	"time"

	"github.com/Coflnet/pr-controller/internal/github"
	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/Coflnet/pr-controller/internal/mongo"
	"github.com/Coflnet/pr-controller/internal/usecase/pr"
	"github.com/rs/zerolog/log"
)

func UpdatePrEnvs() error {

	log.Info().Msgf("updating pr envs")

	githubPrs, err := github.ListPrs()
	if err != nil {
		return err
	}

	dbPrs, err := mongo.ListPrs()
	if err != nil {
		return err
	}

	newRequests := newPrs(githubPrs, dbPrs)
	for _, newPr := range newRequests {
		err = pr.Create(newPr)
		if err != nil {
			log.Error().Err(err).Msgf("there was a problem when creating pr %s", newPr.Repo)
			return err
		}

		log.Info().Msgf("waiting 1 minutes before doing something else")
		time.Sleep(time.Minute * 1)
	}

	updatedPrs := updatedPrs(githubPrs, dbPrs)
	for _, updatePr := range updatedPrs {
		pr.Update(updatePr)
		if err != nil {
			log.Error().Err(err).Msgf("there was a problem when updating pr %s", updatePr.Repo)
			return err
		}

		log.Info().Msgf("waiting 1 minutes before doing something else")
		time.Sleep(time.Minute * 1)
	}

	destroyPrs := destroyPrs(githubPrs, dbPrs)
	for _, destroyPr := range destroyPrs {
		pr.Destroy(destroyPr)
		if err != nil {
			log.Error().Err(err).Msgf("there was a problem when creating pr %s", destroyPr.Repo)
			return err
		}

		log.Info().Msgf("waiting 1 minutes before doing something else")
		time.Sleep(time.Minute * 1)
	}

	log.Info().Msgf("finsihed updating pr envs")

	return nil
}

func newPrs(githubPrs []*model.Pr, dbPrs []*model.Pr) []*model.Pr {
	newPrs := []*model.Pr{}

	for _, githubPr := range githubPrs {

		found := false

		for _, dbPr := range dbPrs {

			if githubPr.Repo == dbPr.Repo &&
				githubPr.Number == dbPr.Number &&
				githubPr.Owner == dbPr.Owner {
				found = true
			}

		}

		if !found {
			newPrs = append(newPrs, githubPr)
		}
	}

	log.Info().Msgf("found %d new prs", len(newPrs))

	return newPrs
}

func updatedPrs(githubPrs []*model.Pr, dbPrs []*model.Pr) []*model.Pr {
	updatePrs := []*model.Pr{}

	for _, githubPr := range githubPrs {
		for _, dbPr := range dbPrs {

			if githubPr.Repo != dbPr.Repo {
				continue
			}

			if githubPr.Number != dbPr.Number {
				continue
			}

			if githubPr.LastCommit == dbPr.LastCommit {
				continue
			}

			updatePrs = append(updatePrs, githubPr)
		}
	}

	log.Info().Msgf("found %d updated prs", len(updatePrs))
	return updatePrs
}

func destroyPrs(githubPrs []*model.Pr, dbPrs []*model.Pr) []*model.Pr {
	destroyPrs := []*model.Pr{}

	for _, dbPr := range dbPrs {
		found := false
		for _, githubPr := range githubPrs {
			if githubPr.Repo == dbPr.Repo &&
				githubPr.Number == dbPr.Number &&
				githubPr.Owner == dbPr.Owner {
				found = true
			}
		}

		if !found {
			destroyPrs = append(destroyPrs, dbPr)
		}
	}

	log.Info().Msgf("found %d destroyed prs", len(destroyPrs))
	return destroyPrs
}
