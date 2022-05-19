package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

type discordWebhookPayload struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

func SendPrCreateMessage(pr *model.Pr) error {
	url := os.Getenv("DISCORD_WEBHOOK")

	payload := &discordWebhookPayload{
		Content:  fmt.Sprintf("new pr env created: %s", pr.CompleteDomain()),
		Username: "GitHub PR Bot",
	}

	j, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msgf("could not serialize discord notification payload")
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		log.Error().Err(err).Msgf("could not create request for discord notification payload")
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("could not send discord notification payload")
		return err
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		log.Error().Msg("failed to send discord webhook notification")
		return err
	}

	log.Info().Msgf("successfully send discord webhook notification")
	return nil
}
