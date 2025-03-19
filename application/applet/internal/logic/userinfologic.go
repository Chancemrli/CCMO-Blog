package logic

import (
	"context"

	"CCMO/gozero/blog/application/applet/internal/code"
	"CCMO/gozero/blog/application/applet/internal/svc"
	"CCMO/gozero/blog/application/applet/internal/types"
	"CCMO/gozero/blog/application/user/rpc/user"
	"CCMO/gozero/blog/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	// userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	// if err != nil {
	// 	return nil, err
	// }
	tokenString := l.ctx.Value("token").(string)
	token, err := jwt.ParseToken(tokenString, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		logx.Infof("ParseToken error: %v", err)
		return nil, err
	}
	userId := token.UserID

	u, err := l.svcCtx.UserRPC.FindById(l.ctx, &user.FindByIdRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("FindById userId: %d error: %v", userId, err)
		return nil, err
	}

	if u.UserId == 0 {
		return nil, code.UserNotExist
	}

	return &types.UserInfoResponse{
		UserId:   u.UserId,
		Username: u.Username,
		Avatar:   u.Avatar,
	}, nil
}
