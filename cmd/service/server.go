package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ovsinc/app-validate-errors-example/internal/service/ports"
	"go.uber.org/fx"
)

func httpServer() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	return app
}

type Routers struct {
	Static fiber.Router
	API    fiber.Router
}

func getGouters(app *fiber.App) Routers {
	apiGroup := app.Group(_api_path)
	apiGroup.Use(func(c *fiber.Ctx) error {
		if c.Is("json") {
			return c.Next()
		}
		return c.Status(http.StatusBadRequest).
			JSON(
				ports.ChangePasswordResponse{
					Common: ports.Common{
						Success: false,
						Message: "Bad request",
					},
					Error: ports.ErrorPayload{
						"common": []string{"Only JSON allowed!"},
					},
				},
			)
	})

	staticGroup := app.Group(_static_path)

	return Routers{
		API:    apiGroup,
		Static: staticGroup,
	}
}

func startService(lifecycle fx.Lifecycle, app *fiber.App) error {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
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
