package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ovsinc/app-validate-errors-example/internal/service/ports"
	"go.uber.org/fx"
)

const (
	_portEnv    = "APP_PORT"
	_port       = 8000
	_heath_path = "/health"
	_api_path   = "/api"
)

func httpServer() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Test API 1.0",
		AppName:       "Test API 1.0",
		//
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})
	app.Use(logger.New())
	return app
}

type Routers struct {
	API     fiber.Router
	Health  fiber.Router
	Monotor fiber.Router
}

func getGouters(app *fiber.App) Routers {
	apiGroup := app.Group(_api_path)
	apiGroup.Use(func(c *fiber.Ctx) error {
		// поддерживаем только JSON
		if !c.Is("json") {
			return c.Status(http.StatusBadRequest).
				JSON(
					ports.ChangePasswordResponse{
						Success: false,
						Error: ports.ErrorPayload{
							"common": []string{"Only JSON allowed!"},
						},
					},
				)
		}

		return c.Next()
	})

	healthGroup := app.Group(_heath_path)

	return Routers{
		API:     apiGroup,
		Health:  healthGroup,
		Monotor: healthGroup,
	}
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
				log.Printf("Stop server on :%d", port)
				return app.Shutdown()
			},
		},
	)
	return nil
}
