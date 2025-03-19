package logic

import (
	"context"

	"CCMO/gozero/blog/application/user/rpc/internal/svc"
	"CCMO/gozero/blog/application/user/rpc/model"
	"CCMO/gozero/blog/application/user/rpc/userService"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	result, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Username: in.Username,
		Password: in.Password,
		Avatar:   in.Avatar,
		Mobile:   in.Mobile,
	})

	if err != nil {
		logx.Errorf("user %s#%s insert error: %v", in.Mobile, in.Username, err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logx.Errorf("user %s#%s get id error: %v", in.Mobile, in.Username, err)
		return nil, err
	}

	return &userService.RegisterResponse{UserId: id}, nil
}
