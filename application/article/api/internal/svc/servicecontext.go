package svc

import (
	"CCMO/gozero/blog/application/article/api/internal/config"
	"CCMO/gozero/blog/application/article/rpc/article"
	"CCMO/gozero/blog/application/egg/rpc/eggclient"
	"CCMO/gozero/blog/application/user/rpc/user"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/zeromicro/go-zero/zrpc"
)

const (
	defaultOssConnectTimeout   = 1
	defaultOssReadWriteTimeout = 3
)

type ServiceContext struct {
	Config     config.Config
	OssClient  *oss.Client
	ArticleRPC article.Article
	UserRPC    user.User
	EggRPC     eggclient.Egg
}

func NewServiceContext(c config.Config) *ServiceContext {
	if c.Oss.ConnectTimeout == 0 {
		c.Oss.ConnectTimeout = defaultOssConnectTimeout
	}
	if c.Oss.ReadWriteTimeout == 0 {
		c.Oss.ReadWriteTimeout = defaultOssReadWriteTimeout
	}
	oc, err := oss.New(c.Oss.Endpoint, c.Oss.AccessKeyId, c.Oss.AccessKeySecret,
		oss.Timeout(c.Oss.ConnectTimeout, c.Oss.ReadWriteTimeout))
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:     c,
		OssClient:  oc,
		ArticleRPC: article.NewArticle(zrpc.MustNewClient(c.ArticleRPC)),
		UserRPC:    user.NewUser(zrpc.MustNewClient(c.UserRPC)),
		EggRPC:     eggclient.NewEgg(zrpc.MustNewClient(c.EggRPC)),
	}
}
