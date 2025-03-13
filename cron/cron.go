package cron

import (
	"bluebell_backend/dao/redis"
	"fmt"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

var job *cron.Cron

func InitCron() {
	job = cron.New()
	// 添加定时任务
	job.AddFunc("*/5 * * * *", func() {
		if err := redis.RefreshLikeRank(); err != nil {
			zap.L().Error("cron RefreshLikeRank failed:", zap.Error(err))
		}
		fmt.Println("Refresh Like Rank")
	})

	job.Start()
}
