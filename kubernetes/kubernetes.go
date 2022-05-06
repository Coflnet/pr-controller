package kubernetes

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

func Init() error {
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Error().Err(err).Msgf("failed to get kubernetes config")
		return err
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create clientset")
		return err
	}

	return nil
}
