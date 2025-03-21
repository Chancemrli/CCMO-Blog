package svc

import (
	"CCMO/gozero/blog/application/egg/rpc/internal/config"
	"CCMO/gozero/blog/pkg/agent/coze"
)

type ServiceContext struct {
	Config config.Config
	Agent  *coze.CozeClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Agent:  coze.NewCozeClient(c.PAT),
	}
}
