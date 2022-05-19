package pr

import (
	metrics "github.com/Coflnet/pr-controller/internal"
	"github.com/Coflnet/pr-controller/internal/kubernetes"
	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/Coflnet/pr-controller/internal/mongo"
	"github.com/rs/zerolog/log"
)

func Destroy(pr *model.Pr) error {

	log.Info().Msgf("destroy pr for branch", pr.Branch)
	err := kubernetes.Destroy(pr)
	if err != nil {
		return err
	}

	log.Info().Msgf("deleteing entry in db for branch %s", pr.Branch)
	err = mongo.DeletePr(pr)
	if err != nil {
		log.Error().Err(err).Msg("error deleting pr from db")
		return err
	}

	log.Info().Msgf("destroyed pr for branch %s", pr.Branch)
	metrics.RemoveEnvironment()

	return nil
}
