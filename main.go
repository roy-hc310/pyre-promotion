package main

import (
	"pyre-promotion/core"
	"pyre-promotion/core-internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	err := utils.LoadGlobalEnv(".")
	if err != nil {
		log.Error().Msg(err.Error())
	}

	g := gin.Default()

	application := core.NewApplication(utils.GlobalEnv)

	err = core.Router(g, application)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	g.Run(":" + utils.GlobalEnv.Port)

}

// migrate create -ext sql -dir ./core/sqlc/migrations -seq init_tables

// migrate -database "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable" -path core-internal/sqlc/migrations up
