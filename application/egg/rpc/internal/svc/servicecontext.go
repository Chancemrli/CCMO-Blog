package svc

import (
	"CCMO/gozero/blog/application/egg/rpc/internal/config"
	"CCMO/gozero/blog/application/egg/rpc/internal/model"
	"CCMO/gozero/blog/pkg/agent/coze"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	Agent      *coze.CozeClient
	ReplyModel model.ReplyModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Agent:      coze.NewCozeClient(c.PAT),
		ReplyModel: model.NewReplyModel(sqlx.NewMysql(c.DataSource)),
	}
}
