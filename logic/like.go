package logic

import (
	"bluebell_backend/dao/redis"
	mykafka "bluebell_backend/kafka"
	"bluebell_backend/models"
	"bluebell_backend/pkg/errcode"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func PostLikeLogic(userID, postID string) (int64, errcode.Status) {
	// 从redis点赞该文章
	likeCount, err := redis.PostLike(userID, postID)
	status := errcode.GetStatus(errcode.AllOK)
	if err != nil && !errors.Is(err, redis.ErrorRedisNotFound) {
		status = errcode.GetStatus(errcode.RedisRunTimeError)
		return 0, status
	}

	uid, _ := strconv.ParseUint(userID, 10, 64)
	pid, _ := strconv.ParseUint(postID, 10, 64)

	// 如果处于rank中，则更新文章热度
	redis.RefreshRankScore(int64(pid))

	// 异步写入kafka消息队列
	userBehavior := models.UserTracking{
		UserID: uid,
		PostID: pid,
		Opt:    1, // 1:like
	}
	// 序列化
	value, err := json.Marshal(&userBehavior)
	if err != nil {
		status = errcode.GetStatus(errcode.ProgramError)
		return 0, status
	}

	if err := mykafka.KWriter.WriteMessages(context.Background(), kafka.Message{
		Value: value,
	}); err != nil {
		status = errcode.GetStatus(errcode.ProgramError)
		return 0, status
	}

	return likeCount, status
}

func PostUnlikeLogic(userID, postID string) (int64, errcode.Status) {
	// 从redis点踩该文章
	unlikeCount, err := redis.PostUnlike(userID, postID)
	status := errcode.GetStatus(errcode.AllOK)
	if err != nil && !errors.Is(err, redis.ErrorRedisNotFound) {
		status = errcode.GetStatus(errcode.RedisRunTimeError)
		return 0, status
	}

	uid, _ := strconv.ParseUint(userID, 10, 64)
	pid, _ := strconv.ParseUint(postID, 10, 64)

	// 如果处于rank中，则更新文章热度
	redis.RefreshRankScore(int64(pid))

	// 异步写入kafka消息队列
	userBehavior := models.UserTracking{
		UserID: uid,
		PostID: pid,
		Opt:    -1, // 1:like
	}
	// 序列化
	value, err := json.Marshal(&userBehavior)
	if err != nil {
		status = errcode.GetStatus(errcode.ProgramError)
		return 0, status
	}

	if err := mykafka.KWriter.WriteMessages(context.Background(), kafka.Message{
		Value: value,
	}); err != nil {
		status = errcode.GetStatus(errcode.ProgramError)
		return 0, status
	}

	return unlikeCount, status
}
