package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/valyala/fasthttp/expvarhandler"
	"go.uber.org/fx"

	"github.com/ovsinc/app-validate-errors-example/internal/service/ports"
)

const (
	_varsPath = "/debug/vars"
)

func registryAPIHandler(
	lifecycle fx.Lifecycle,
	routers Routers,
	h ports.PasswordChange,
	bundle *i18n.Bundle,
) error {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				hdlr := &handler{h: h, bundle: bundle}
				routers.API.Add(http.MethodPost, "/v1", hdlr.ChangePassword)
				routers.Health.Add(http.MethodGet, "/", hdlr.Health)
				routers.Monotor.Add(http.MethodGet, _varsPath, func(c *fiber.Ctx) error {
					expvarhandler.ExpvarHandler(c.Context())
					return nil
				})
				return nil
			},
		},
	)
	return nil
}

type handler struct {
	h      ports.PasswordChange
	bundle *i18n.Bundle
}

const (
	ParseRequestFailed = "ParseRequestFailed"
)

var MsgBadRequest = i18n.Message{
	ID:          ParseRequestFailed,
	Description: "Ошибка анализа запроса",
	Other:       "Request analysis error",
}

func (h *handler) ChangePassword(c *fiber.Ctx) error {
	accept := c.Request().Header.Peek("Accept-Language")
	localizer := i18n.NewLocalizer(h.bundle, string(accept))

	req := new(ports.ChangePasswordRequest)
	if err := c.BodyParser(&req); err != nil {
		msg := MsgBadRequest.Other
		if lmsg, lerr := localizer.LocalizeMessage(&MsgBadRequest); lerr == nil {
			msg = lmsg
		}
		return c.Status(http.StatusBadRequest).
			JSON(
				ports.ChangePasswordResponse{
					Success: false,
					Error: ports.ErrorPayload{
						"common": []string{msg},
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

func (h *handler) Health(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
