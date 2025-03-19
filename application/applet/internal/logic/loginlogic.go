package logic

import (
	"context"
	"strings"

	"CCMO/gozero/blog/application/applet/internal/code"
	"CCMO/gozero/blog/application/applet/internal/svc"
	"CCMO/gozero/blog/application/applet/internal/types"
	"CCMO/gozero/blog/application/user/rpc/userService"
	"CCMO/gozero/blog/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 规范化数据格式
	req.Mobile = strings.TrimSpace(req.Mobile)
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)

	// 验证手机账号是否存在
	user, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &userService.FindByMobileRequest{Mobile: req.Mobile})
	if err != nil {
		logx.Errorf("FindByMobile error: %v", err)
		return nil, err
	}
	// 手机号不存在
	if user.UserId == 0 || user == nil {
		logx.Infof("User not found. mobile: %s", req.Mobile)
		return nil, code.LoginMobileInvalid
	}

	// 对比验证码
	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Infof("checkVerificationCode error: %v", err)
		return nil, err
	}

	// 生成token
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": user.UserId,
		},
	})
	if err != nil {
		logx.Errorf("BuildTokens error: %v", err)
		return nil, err
	}

	// 删除缓存中的验证码
	_ = delActivationCache(req.Mobile, l.svcCtx.BizRedis)

	return &types.LoginResponse{
		UserId: user.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}
