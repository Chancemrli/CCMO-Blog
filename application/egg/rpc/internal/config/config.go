package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	User struct {
		UserID   uint64
		Username string
		Password string
	}
	PAT        string
	WID        string
	InPut      string
	DataSource string
}
