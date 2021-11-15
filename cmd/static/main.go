package main

import (
	"log"
	"time"

	"go.uber.org/fx"
)

const (
	GracefulStopTimeout = 10 * time.Second
	StartTimeout        = 10 * time.Second
)

func main() {
	appCtx := fx.New(
		// options
		fx.StartTimeout(StartTimeout),
		fx.StopTimeout(GracefulStopTimeout),

		fx.Provide(
			httpServer,
		),

		fx.Invoke(
			startService,
		),
	)

	appCtx.Run()

	if err := appCtx.Err(); err != nil {
		log.Println(err)
	}
}
