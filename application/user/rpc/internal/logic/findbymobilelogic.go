package logic

import (
	"context"
	"errors"

	"CCMO/gozero/blog/application/user/rpc/internal/svc"
	"CCMO/gozero/blog/application/user/rpc/userService"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FindByMobileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByMobileLogic {
	return &FindByMobileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByMobileLogic) FindByMobile(in *userService.FindByMobileRequest) (*userService.FindByMobileResponse, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil && !errors.Is(err, sqlx.ErrNotFound) {
		logx.Errorf("FindByMobile mobile: %s error: %v", in.Mobile, err)
		return nil, err
	}
	if user == nil {
		return &userService.FindByMobileResponse{}, nil
	}

	return &userService.FindByMobileResponse{
		UserId:   int64(user.Id),
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
