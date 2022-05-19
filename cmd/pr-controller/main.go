package main

import (
	"github.com/Coflnet/pr-controller/internal/api"
	"github.com/Coflnet/pr-controller/internal/github"
	"github.com/Coflnet/pr-controller/internal/kubernetes"
	"github.com/Coflnet/pr-controller/internal/metrics"
	"github.com/Coflnet/pr-controller/internal/mongo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {

	go github.Init()

	go func() {
		err := kubernetes.Init()
		if err != nil {
			log.Fatal().Err(err).Msgf("could not initialize kubernetes")
		}
	}()

	go func() {
		err := mongo.Init()
		if err != nil {
			log.Fatal().Err(err).Msg("could not connect to mongo")
		}
	}()
	defer mongo.Disconnect()

	go func() {
		err := metrics.Init()
		if err != nil {
			log.Fatal().Err(err).Msg("could not connect to metrics")
		}
	}()

	err := api.StartApi()
	log.Error().Err(err).Msgf("api stopped")
}
