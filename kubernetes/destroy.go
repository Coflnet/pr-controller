package kubernetes

import (
	"context"

	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Destroy(pr *model.Pr) error {
	log.Warn().Msgf("destroying not implmented yet, destroy pr %s/%s", pr.Owner, pr.Repo)
	return nil
}

func DestroyDeployment(pr *model.Pr) error {

	deploymentsClient := clientset.AppsV1().Deployments("pr-env")

	err := deploymentsClient.Delete(context.TODO(), pr.KubernetesResourceName(), metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
