package usecase

import (
	"github.com/Coflnet/pr-controller/internal/metrics"
	"sync"
	"time"

	"github.com/Coflnet/pr-controller/internal/github"
	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/Coflnet/pr-controller/internal/mongo"
	"github.com/Coflnet/pr-controller/internal/usecase/pr"
	"github.com/rs/zerolog/log"
)

var wg sync.WaitGroup

func UpdatePrEnvs() error {

	wg.Wait()
	wg.Add(1)
	defer wg.Done()

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
			metrics.AddError()
			return err
		}

		log.Info().Msgf("waiting 10 minutes before doing something else")
		time.Sleep(time.Minute * 10)
	}

	updatedPrs := updatedPrs(githubPrs, dbPrs)
	for _, updatePr := range updatedPrs {
		err = pr.Update(updatePr)
		if err != nil {
			log.Error().Err(err).Msgf("there was a problem when updating pr %s", updatePr.Repo)
			metrics.AddError()
			return err
		}

		log.Info().Msgf("waiting 10 minutes before doing something else")
		time.Sleep(time.Minute * 10)
	}

	destroyPrs := destroyPrs(githubPrs, dbPrs)
	for _, destroyPr := range destroyPrs {
		err = pr.Destroy(destroyPr)
		if err != nil {
			log.Error().Err(err).Msgf("there was a problem when creating pr %s", destroyPr.Repo)
			metrics.AddError()
			return err
		}

		log.Info().Msgf("waiting 10 minutes before doing something else")
		time.Sleep(time.Minute * 10)
	}

	log.Info().Msgf("finished updating pr envs")

	return nil
}

func newPrs(githubPrs []*model.Pr, dbPrs []*model.Pr) []*model.Pr {
	var newPrs []*model.Pr

	for _, githubPr := range githubPrs {

		found := false

		for _, dbPr := range dbPrs {

			log.Info().Msgf("comparing github pr: %s, %s, %d with mongo pr %s, %s, %d", githubPr.Owner, githubPr.Repo, githubPr.Number, dbPr.Owner, dbPr.Repo, dbPr.Number)

			if githubPr.Repo == dbPr.Repo &&
				githubPr.Number == dbPr.Number &&
				githubPr.Owner == dbPr.Owner {
				found = true
				log.Info().Msgf("comparison above returns true")
				break
			}

			log.Info().Msgf("comparison above returns false")
		}

		if !found {
			newPrs = append(newPrs, githubPr)
		}
	}

	log.Info().Msgf("found %d new prs", len(newPrs))

	return newPrs
}

func updatedPrs(githubPrs []*model.Pr, dbPrs []*model.Pr) []*model.Pr {
	var updatePrs []*model.Pr

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
	var destroyPrs []*model.Pr

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
