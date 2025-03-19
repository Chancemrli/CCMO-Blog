package logic

import (
	"context"

	"CCMO/gozero/blog/application/user/rpc/internal/svc"
	"CCMO/gozero/blog/application/user/rpc/userService"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsLogic {
	return &SendSmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSmsLogic) SendSms(in *userService.SendSmsRequest) (*userService.SendSmsResponse, error) {
	// todo: add your logic here and delete this line

	return &userService.SendSmsResponse{}, nil
}
