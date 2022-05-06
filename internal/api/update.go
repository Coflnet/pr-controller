package api

import (
	"github.com/Coflnet/pr-controller/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func triggerUpdate(c *gin.Context) {
	log.Info().Msgf("got a update request")

	go usecase.UpdatePrEnvs()

	log.Info().Msgf("update was successful")
}
