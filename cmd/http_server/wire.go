//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/jr-dragon/dynamic_link/internal/biz"

	"github.com/jr-dragon/dynamic_link/api"
	"github.com/jr-dragon/dynamic_link/internal/data"
	"github.com/jr-dragon/dynamic_link/internal/server"
)

func wireApp(cfg data.Config) (*App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, api.ProviderSet, biz.ProviderSet, newApp))
}
