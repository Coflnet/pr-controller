package main

import (
	"github.com/Coflnet/pr-controller/internal/api"
	"github.com/Coflnet/pr-controller/internal/github"
	"github.com/Coflnet/pr-controller/internal/kubernetes"
	"github.com/Coflnet/pr-controller/internal/mongo"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {

	github.Init()

	kubernetes.Init()

	err := mongo.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to mongo")
	}
	defer mongo.Disconnect()

	err = api.StartApi()
	log.Error().Err(err).Msgf("api stopped")
}
