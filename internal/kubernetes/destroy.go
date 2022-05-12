package kubernetes

import (
	"context"

	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Destroy(pr *model.Pr) error {
	err := DestroyDeployment(pr)
	if err != nil {
		return err
	}

	err = DestroyService(pr)
	if err != nil {
		return err
	}

	err = DestroyIngress(pr)
	return err
}

func DestroyDeployment(pr *model.Pr) error {

	deploymentsClient := clientset.AppsV1().Deployments("pr-env")

	err := deploymentsClient.Delete(context.TODO(), pr.KubernetesResourceName(), metav1.DeleteOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("could not delete deployment %s", pr.KubernetesResourceName())
		return err
	}

	return nil
}

func DestroyService(pr *model.Pr) error {
	serviceClient := clientset.CoreV1().Services("pr-env")

	err := serviceClient.Delete(context.TODO(), pr.KubernetesResourceName(), metav1.DeleteOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("could not delete service %s", pr.KubernetesResourceName())
		return err
	}

	return nil
}

func DestroyIngress(pr *model.Pr) error {
	ingressClient := clientset.NetworkingV1().Ingresses("pr-env")

	err := ingressClient.Delete(context.TODO(), pr.KubernetesResourceName(), metav1.DeleteOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("error deleting ingress %s", pr.KubernetesResourceName())
		return err
	}

	return nil
}
