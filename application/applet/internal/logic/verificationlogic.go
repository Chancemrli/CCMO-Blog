package logic

import (
	"context"
	"fmt"
	"strconv"

	"CCMO/gozero/blog/application/applet/internal/svc"
	"CCMO/gozero/blog/application/applet/internal/types"
	"CCMO/gozero/blog/pkg/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	prefixActivation        = "biz#activation#%s"
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	// todo: add your logic here and delete this line
	// 检查该手机号获取验证码的次数
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("getVerificationCount mobile: %s error: %v", req.Mobile, err)
	}
	// 验证次数超出限制，拒绝访问
	if count >= verificationLimitPerDay {
		return nil, fmt.Errorf("daily verification limit reached")
	}

	// 获取有效验证码
	code, err := getActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("getActivationCache mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	// 当前没有有效的验证码，则生成六位数的随机验证码
	if len(code) == 0 {
		code = util.RandomNumeric(6)
	}

	// // 不知道干嘛的
	// _, err = l.svcCtx.UserRPC.SendSms(l.ctx, &user.SendSmsRequest{
	// 	Mobile: req.Mobile,
	// })
	// if err != nil {
	// 	logx.Errorf("sendSms mobile: %s error: %v", req.Mobile, err)
	// 	return nil, err
	// }

	// 上传验证码
	err = saveActivetionCache(req.Mobile, code, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("saveActivetionCache mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	// 增加手机号短信验证次数
	_, err = incrVerificationCount(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
	}

	return &types.VerificationResponse{}, nil
}

// 获取当前手机号使用短信验证次数
func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, nil
	}

	return strconv.Atoi(val)
}

// 获取有效验证码
func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Get(key)
}

// 生成验证码存入redis
func saveActivetionCache(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}

// 增加短信验证次数
func incrVerificationCount(mobile string, rds *redis.Redis) (int64, error) {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	rds.SetnxEx(key, "0", 86400)
	n, err := rds.Incr(key)
	return n, err
}

// 删除缓存中的验证码
func delActivationCache(mobile string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	_, err := rds.Del(key)
	return err
}
