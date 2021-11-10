package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"

	"github.com/ovsinc/app-validate-errors-example/internal/service/ports"
)

func registryAPIHandler(lifecycle fx.Lifecycle, routers Routers, h ports.PasswordChange) error {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				hdlr := handler{h: h}
				routers.API.Add(http.MethodPost, "/v1", hdlr.ChangePassword)
				return nil
			},
		},
	)
	return nil
}

type handler struct {
	h ports.PasswordChange
}

func (h *handler) ChangePassword(c *fiber.Ctx) error {
	req := new(ports.ChangePasswordRequest)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(
				ports.ChangePasswordResponse{
					Common: ports.Common{
						Success: false,
						Message: "Bad request",
					},
					Error: ports.ErrorPayload{
						"common": &ports.ErrorOut{Message: "Request is bad"},
					},
				},
			)
	}

	ctx := context.Background()

	resp, code, err := h.h.ChangePassword(ctx, req)
	if err != nil {
		log.Println(err)
	}

	return c.Status(code).JSON(resp)
}
