package main

import (
	"log"
	"time"

	"go.uber.org/fx"

	"github.com/ovsinc/app-validate-errors-example/internal/service/adapter"
	"github.com/ovsinc/app-validate-errors-example/internal/service/app"
	"github.com/ovsinc/app-validate-errors-example/internal/service/ports"
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

			ports.NewHttpServer,
			app.NewCheckPassword,
			app.NewChangePassword,
			adapter.NewPasswordRepository,
		),

		fx.Invoke(
			registryStaticHandler,
			startService,
		),
	)

	appCtx.Run()

	if err := appCtx.Err(); err != nil {
		log.Println(err)
	}
}
