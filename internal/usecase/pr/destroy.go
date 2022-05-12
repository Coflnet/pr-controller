package pr

import (
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
	mongo.DeletePr(pr)

	return nil
}
