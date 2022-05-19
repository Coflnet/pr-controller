package argo

import (
	"bytes"
	"encoding/json"
	metrics "github.com/Coflnet/pr-controller/internal"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

type WebhookRequest struct {
	Image   string `json:"image"`
	Tag     string `json:"tag"`
	Branch  string `json:"branch"`
	Git_url string `json:"git_url"`
}

func BuildImage(image, tag, repo, branch string) error {

	url := os.Getenv("BUILD_IMAGE_WEBHOOK")

	webhookRequest := WebhookRequest{
		Image:   image,
		Tag:     tag,
		Git_url: repo,
		Branch:  branch,
	}

	log.Debug().Msgf("sending build request with data: %+v", webhookRequest)

	j, err := json.Marshal(webhookRequest)
	if err != nil {
		log.Error().Err(err).Msgf("could not serialize webhook request")
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		log.Error().Err(err).Msgf("could not create request for %s", url)
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Error().Err(err).Msgf("there was an error when triggering the build of the image %s:%s", image, tag)
		return err
	}

	if response.StatusCode != 200 {
		log.Error().Err(err).Msgf("there was an error when triggering the build of the image %s:%s", image, tag)
		return err
	}

	log.Info().Msgf("triggered a build of container %s:%s, git repo: %s", image, tag, repo)
	metrics.CIPipelineTriggered()

	return nil
}
