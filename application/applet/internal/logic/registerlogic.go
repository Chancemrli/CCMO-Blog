package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"CCMO/gozero/blog/application/applet/internal/code"
	"CCMO/gozero/blog/application/applet/internal/svc"
	"CCMO/gozero/blog/application/applet/internal/types"
	"CCMO/gozero/blog/application/user/rpc/userService"
	"CCMO/gozero/blog/pkg/encrypt"
	"CCMO/gozero/blog/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// todo: add your logic here and delete this line
	// 规范输入
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)
	req.Mobile = strings.TrimSpace(req.Mobile)

	// 参数校验
	if len(req.Mobile) == 0 {
		return nil, code.RegisterMobileEmpty
	}

	// 验证码校验
	err = checkVerificationCode(l.svcCtx.BizRedis, req.Mobile, req.VerificationCode)
	if err != nil {
		logx.Errorf("checkVerificationCode error: %v", err)
		return nil, err
	}

	// // 加密手机号
	// mobile, err := encrypt.EncMobile(req.Mobile)
	// fmt.Println(mobile)
	// if err != nil {
	// 	logx.Errorf("EncMobile mobile: %s error: %v", req.Mobile, err)
	// 	return nil, err
	// }

	// 加密密码
	password := encrypt.EncPassword(req.Password)
	fmt.Println(password)
	if err != nil {
		logx.Errorf("EncPassword password: %s error: %v", req.Password, err)
		return nil, err
	}

	// 调用rpc服务查看手机号是否已注册
	user, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &userService.FindByMobileRequest{
		Mobile: req.Mobile,
	})
	// 返回错误
	if err != nil {
		logx.Errorf("FindByMobile error: %v", err)
		return nil, err
	}
	// 手机号已经被创建
	if user.UserId > 0 {
		logx.Info("Mobie has been registered. Please use another mobile number.  mobile: " + req.Mobile)
		return nil, code.MobileHasRegistered
	}

	// TODO调用rpc服务创建用户
	uid, err := l.svcCtx.UserRPC.Register(l.ctx, &userService.RegisterRequest{
		Username: req.Name,
		Password: password,
		Avatar:   "",
		Mobile:   req.Mobile,
	})
	if err != nil {
		logx.Errorf("Register error: %v", err)
		return nil, err
	}

	// 创建token
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.AccessExpire,
		Fields: map[string]interface{}{
			"userId": uid.UserId,
		},
	})

	if err != nil {
		logx.Errorf("BuildTokens error: %v", err)
		return nil, err
	}

	_ = delActivationCache(req.Mobile, l.svcCtx.BizRedis)

	return &types.RegisterResponse{
		UserId: uid.UserId,
		Token: types.Token{
			AccessToken:  token.AccessToken,
			AccessExpire: token.AccessExpire,
		},
	}, nil
}

func checkVerificationCode(rds *redis.Redis, mobile, code string) error {
	cacheCode, err := getActivationCache(mobile, rds)
	if err != nil {
		return err
	}
	if cacheCode == "" {
		return errors.New("verification code expired")
	}
	if cacheCode != code {
		return errors.New("verification code failed")
	}

	return nil
}
