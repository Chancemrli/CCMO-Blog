package logic

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"bluebell_backend/models"
	"bluebell_backend/pkg/errcode"
	"errors"
	"strconv"
)

const (
	TopNum    = 20
	BackupNum = 10
)

func RankByLikeCountLogic() (rank []models.PostInScoreRank, status errcode.Status) {
	// 从Redis获取排行榜
	rank, err := redis.GetRankByScoreCount(TopNum)
	status = errcode.GetStatus(errcode.AllOK)
	if err != nil && !errors.Is(err, redis.ErrorRedisNotFound) {
		status = errcode.GetStatus(errcode.RedisRunTimeError)
		return
	}

	// redis缓存失效
	if rank == nil || len(rank) < TopNum {
		// 从数据库获取前TopNum+BackNum的排行
		rank, err = mysql.GetRankByScore(TopNum + BackupNum)
		if err != nil {
			status = errcode.GetStatus(errcode.SQLRunTimeError)
			return
		}

		// 存入redis
		if err := redis.InitRank(rank); err != nil {
			status = errcode.GetStatus(errcode.RedisRunTimeError)
			return
		}

		// 将对应文章缓存为热点
		for _, post := range rank {
			pid := strconv.Itoa(int(post.PostID))
			GetPost(pid)
		}

		// 取出前TopNum的排行
		rank = rank[:TopNum]
	}

	return
}
