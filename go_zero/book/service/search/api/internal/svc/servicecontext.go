package svc

import (
	"book/pkg/interceptor/rpcclient"
	"book/service/search/api/internal/config"
	"book/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userclient.NewUser(zrpc.MustNewClient(c.UserRpc,zrpc.WithUnaryClientInterceptor(rpcclient.BreakerInterceptor))),
	}
}
