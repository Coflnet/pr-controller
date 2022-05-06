package pr

import (
	"github.com/Coflnet/pr-controller/internal/model"
)

func Update(pr *model.Pr) error {

	err := Destroy(pr)
	if err != nil {
		return err
	}

	err = Create(pr)
	if err != nil {
		return err
	}

	return nil
}
