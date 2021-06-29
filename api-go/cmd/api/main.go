package main

import (
	"github.com/Hongbo-Miao/hongbomiao.com/api-go/internal/routes"
	"github.com/Hongbo-Miao/hongbomiao.com/api-go/internal/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	utils.InitLogger()
	var config = utils.GetConfig()
	log.Info().Str("env", config.Env).Str("port", config.Port).Msg("main")

	r := routes.SetupRouter()
	_ = r.Run(":" + config.Port)
}
