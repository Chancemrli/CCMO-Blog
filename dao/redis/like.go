package redis

import (
	"bluebell_backend/pkg/util"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

func PostLike(userID, postID string) (int64, error) {
	key := KeyHotArticlePrefix + postID
	exist, err := client.Exists(key).Result()
	// 程序出错
	if err != nil {
		err = ErrorRedisRunTimeError
		return 0, err
	}
	// 键不存在
	if exist == 0 {
		err = ErrorRedisNotFound
		return 0, err
	}
	// 点赞数+1
	res, err := client.HIncrBy(key, "like_count", 1).Result()
	if err != nil {
		err = ErrorRedisRunTimeError
		return 0, err
	}

	return res, nil
}

func PostUnlike(userID, postID string) (int64, error) {
	key := KeyHotArticlePrefix + postID
	exist, err := client.Exists(key).Result()
	// 程序出错
	if err != nil {
		err = ErrorRedisRunTimeError
		return 0, err
	}
	// 键不存在
	if exist == 0 {
		err = ErrorRedisNotFound
		return 0, err
	}
	// 点踩数+1
	res, err := client.HIncrBy(key, "unlike_count", 1).Result()
	if err != nil {
		err = ErrorRedisRunTimeError
		return 0, err
	}

	return res, nil
}

// 更新rank中的分值
func RefreshRankScore(postID int64) error {
	// 判断是否在rank中
	ok, err := IsInRank(postID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrorRedisNotFound
	}

	// 更新文章热度
	detail, err := GetArticleByID(strconv.Itoa(int(postID)))
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%d%s%s", postID, SplitKey, detail.Title)
	score := util.Hot(int(detail.LikeCount), int(detail.UnlikeCount), detail.CreateTime)
	_, err = client.ZAdd(RankByLikeCountKey, redis.Z{
		Member: key,
		Score:  float64(score),
	}).Result()

	if err != nil {
		return ErrorRedisRunTimeError
	}

	return nil
}
