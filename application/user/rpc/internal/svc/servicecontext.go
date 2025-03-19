package svc

import (
	"CCMO/gozero/blog/application/user/rpc/internal/config"
	"CCMO/gozero/blog/application/user/rpc/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlx.NewSqlConn("mysql", c.DataSource), c.CacheRedis),
	}
}
