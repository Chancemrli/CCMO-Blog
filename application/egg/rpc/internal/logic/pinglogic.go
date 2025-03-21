package logic

import (
	"context"

	"CCMO/gozero/blog/application/egg/rpc/egg"
	"CCMO/gozero/blog/application/egg/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *egg.EggRequest) (*egg.EggResponse, error) {
	// todo: add your logic here and delete this line

	return &egg.EggResponse{}, nil
}
