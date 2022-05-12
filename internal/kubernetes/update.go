package kubernetes

import (
	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/rs/zerolog/log"
)

func Update(pr *model.Pr) error {
	log.Warn().Msgf("updating of kubernetes resources not yet implemented")
	return nil
}

func UpdateDeployment(pr *model.Pr) error {
	err := Destroy(pr)
	if err != nil {
		return err
	}

	err = Create(pr)
	return err
}
