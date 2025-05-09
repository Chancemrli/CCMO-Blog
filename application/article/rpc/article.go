package main

import (
	"flag"
	"fmt"

	"CCMO/gozero/blog/application/article/rpc/internal/config"
	"CCMO/gozero/blog/application/article/rpc/internal/server"
	"CCMO/gozero/blog/application/article/rpc/internal/svc"
	"CCMO/gozero/blog/application/article/rpc/pb"
	"CCMO/gozero/blog/pkg/consul"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.DisableStat()
	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterArticleServer(grpcServer, server.NewArticleServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 服务注册
	err := consul.Register(c.Consul, fmt.Sprintf("%s:%d", c.ServiceConf.Prometheus.Host, c.ServiceConf.Prometheus.Port))
	if err != nil {
		fmt.Printf("register consul error: %v\n", err)
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
