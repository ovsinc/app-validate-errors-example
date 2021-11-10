package main

import (
	"context"
	"path"

	"go.uber.org/fx"
)

func registryStaticHandler(lifecycle fx.Lifecycle, path string, routers Routers) error {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				routers.Static.Static(_static_path, path)
				return nil
			},
		},
	)
	return nil
}

func getStaticDir() string {
	return path.Clean(staticDir)
}
