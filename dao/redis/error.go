package redis

import "errors"

var (
	ErrorVoteTimeExpire = errors.New("已过投票时间")
	ErrorVoted          = errors.New("已经投过票了")

	ErrorRedisNotFound     = errors.New("缓存无效")
	ErrorRedisRunTimeError = errors.New("Redis运行时错误")
)
