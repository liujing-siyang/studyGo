//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"realworld/internal/biz"
	"realworld/internal/conf"
	"realworld/internal/data"
	"realworld/internal/server"
	"realworld/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *conf.Jwt, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSetServer, data.ProviderSetData, biz.ProviderSetBiz, service.ProviderSetService, newApp))
}
