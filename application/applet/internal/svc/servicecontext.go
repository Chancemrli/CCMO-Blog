package svc

import (
	"CCMO/gozero/blog/application/applet/internal/config"
	"CCMO/gozero/blog/application/user/rpc/user"
	"CCMO/gozero/blog/pkg/interceptors"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config   config.Config
	UserRPC  user.User
	BizRedis *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	userRPC := zrpc.MustNewClient(c.UserRPC, zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor()))
	return &ServiceContext{
		Config:   c,
		UserRPC:  user.NewUser(userRPC),
		BizRedis: redis.MustNewRedis(c.BizRedis),
	}
}
