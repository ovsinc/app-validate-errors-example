package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	expvarmw "github.com/gofiber/fiber/v2/middleware/expvar"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/fx"
)

const (
	_static_path = "/"
	_portEnv     = "APP_PORT"
	_port        = 3000
)

// Embed a directory
//go:embed dist/spa
var embedDirStatic embed.FS

func httpServer() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Test Static 1.0",
		AppName:       "Test Static 1.0",
		//
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,

		EnableTrustedProxyCheck: false,
		ProxyHeader:             "X-Real-Ip",
	})
	app.Use(logger.New())
	app.Use(
		_static_path,
		filesystem.New(
			filesystem.Config{
				Root:       http.FS(embedDirStatic),
				PathPrefix: "dist/spa",
				Browse:     true,
				Index:      "index.html",
				MaxAge:     600,
			},
		),
	)
	app.Use(expvarmw.New())
	return app
}

func startService(lifecycle fx.Lifecycle, app *fiber.App) error {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port, err := strconv.Atoi(os.Getenv(_portEnv))
				if err != nil {
					log.Printf("port error: %v", err)
					log.Printf("use default port: %d", _port)
					port = _port
				}

				log.Printf("HTTP server listen :%d", port)
				go func() {
					_ = app.Listen(fmt.Sprintf(":%d", port))
				}()
				return nil
			},
			OnStop: func(context.Context) error {
				log.Println("Stop server on")
				return app.Shutdown()
			},
		},
	)
	return nil
}
