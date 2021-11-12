package main

import (
	"flag"
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
	_static_path        = "/"
)

var (
	port      int
	staticDir string
)

func main() {
	flag.IntVar(&port, "port", 8000, "The port to listen on")
	flag.StringVar(&staticDir, "static", "./dist/spa/", "Static directory path")

	flag.Parse()

	appCtx := fx.New(
		// options
		fx.StartTimeout(StartTimeout),
		fx.StopTimeout(GracefulStopTimeout),

		fx.Provide(
			httpServer,

			getGouters,

			ports.NewHttpServer,
			app.NewCheckPassword,
			app.NewChangePassword,
			adapter.NewPasswordRepository,
		),

		fx.Invoke(
			registryAPIHandler,
			startService,
		),
	)

	appCtx.Run()

	if err := appCtx.Err(); err != nil {
		log.Println(err)
	}
}
