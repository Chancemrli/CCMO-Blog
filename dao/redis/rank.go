package redis

import (
	"bluebell_backend/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

const (
	RankByLikeCountKey = "cache:rank:likecount"
	SplitKey           = "{#}"
)

func RefreshLikeRank() error {
	_, err := client.Del(RankByLikeCountKey).Result()
	if err != nil {
		return err
	}
	return nil
}

func InitRank(posts []models.PostInScoreRank) error {
	for _, post := range posts {
		_, err := client.ZAdd(RankByLikeCountKey, redis.Z{
			Member: fmt.Sprintf("%d%s%s", post.PostID, SplitKey, post.Title),
			Score:  float64(post.Score),
		}).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRankByScoreCount(topNum int64) ([]models.PostInScoreRank, error) {
	list, err := client.ZRevRangeWithScores(RankByLikeCountKey, 0, topNum-1).Result()
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, ErrorRedisNotFound
	}

	posts := make([]models.PostInScoreRank, 0)
	for _, z := range list {
		idAndTitle := strings.SplitN(z.Member.(string), SplitKey, 2)
		postID, _ := strconv.ParseUint(idAndTitle[0], 10, 64)
		posts = append(posts, models.PostInScoreRank{
			PostID: postID,
			Title:  idAndTitle[1],
			Score:  int64(z.Score),
		})
	}
	return posts, nil
}

// 判断文章是否在rank中
func IsInRank(postID int64) (bool, error) {
	keys, err := client.ZRange(RankByLikeCountKey, 0, -1).Result()
	if err != nil {
		err = ErrorRedisRunTimeError
		return false, err
	}
	// 遍历比对postID
	for _, key := range keys {
		idAndTitle := strings.SplitN(key, SplitKey, 2)
		postIDStr := idAndTitle[0]
		id, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			return false, err
		}
		if id == postID {
			return true, nil
		}
	}
	return false, nil
}

// 获取rank中文章的排名
