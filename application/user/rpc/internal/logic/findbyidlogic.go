package logic

import (
	"context"
	"errors"

	"CCMO/gozero/blog/application/user/rpc/internal/svc"
	"CCMO/gozero/blog/application/user/rpc/userService"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type FindByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindByIdLogic {
	return &FindByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindByIdLogic) FindById(in *userService.FindByIdRequest) (*userService.FindByIdResponse, error) {
	// todo: add your logic here and delete this line
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(in.UserId))
	if errors.Is(err, sqlx.ErrNotFound) {
		return &userService.FindByIdResponse{}, nil
	} else if err != nil {
		logx.Errorf("rpc FindOne userId: %d error: %v", in.UserId, err)
		return nil, err
	}
	return &userService.FindByIdResponse{
		UserId:   int64(user.Id),
		Username: user.Username,
		Mobile:   user.Mobile,
		Avatar:   user.Avatar,
	}, nil
}
