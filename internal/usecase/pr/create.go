package pr

import (
	"github.com/Coflnet/pr-controller/internal/argo"
	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/Coflnet/pr-controller/internal/mongo"
	"github.com/Coflnet/pr-controller/kubernetes"
	"github.com/rs/zerolog/log"
)

func Create(pr *model.Pr) error {

	// trigger the build of a container
	err := argo.BuildImage(pr.Image, pr.Tag(), pr.GitUrl(), pr.Branch)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when building the container image")
		return err
	}
	log.Info().Msgf("triggered a build for image %s", pr.Image)

	// create the kuberentes deployment
	err = kubernetes.Create(pr)
	if err != nil {
		log.Error().Err(err).Msgf("there was an problem when creating kuberntes resources for pr %s/%s", pr.Owner, pr.Repo)
		return err
	}
	log.Info().Msgf("created kubernetes resources for repo %s", pr.Repo)

	// save the new state in the database
	err = mongo.InsertPr(pr)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error when saving the pr %s/%s", pr.Owner, pr.Repo)
		return err
	}

	log.Warn().Msgf("creating not implmented yet, create pr %s", pr.Repo)
	return nil
}
