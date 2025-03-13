package redis

import (
	"bluebell_backend/models"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func IsHotArticle(id string) (string, error) {
	key, err := client.HGet(KeyHotArticlesMap, id).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	if errors.Is(err, redis.Nil) {
		return "", nil
	}

	return key, nil
}

func GetArticleByID(id string) (post *models.ApiPostDetail, err error) {
	key := KeyHotArticlePrefix + id
	value, err := client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	if value["post_id"] == "0" {
		return nil, errors.New("empty article")
	}

	postID, _ := strconv.ParseUint(id, 10, 64)
	authorID, _ := strconv.ParseUint(value["author_id"], 10, 64)
	communityID, _ := strconv.ParseInt(value["community_id"], 10, 64)
	status, _ := strconv.ParseInt(value["status"], 10, 32)
	createTime, _ := time.Parse("2006-01-02 15:04:05", value["create_time"])
	likeCount, _ := strconv.ParseInt(value["like_count"], 10, 64)
	unlikeCount, _ := strconv.ParseInt(value["unlike_count"], 10, 64)
	post = &models.ApiPostDetail{
		Post: &models.Post{
			PostID:      postID,
			Title:       value["title"],
			Content:     value["content"],
			AuthorId:    authorID,
			CommunityID: communityID,
			Status:      int32(status),
			CreateTime:  createTime,
			LikeCount:   likeCount,
			UnlikeCount: unlikeCount,
		},
		AuthorName:    value["author_name"],
		CommunityName: value["community_name"],
	}
	return post, nil
}

func SetHotArticle(post models.ApiPostDetail) (string, error) {
	// 获取键
	key := fmt.Sprintf("%s%d", KeyHotArticlePrefix, post.PostID)
	value := make(map[string]interface{})
	value["post_id"] = post.PostID
	value["title"] = post.Title
	value["content"] = post.Content
	value["author_id"] = post.AuthorId
	value["author_name"] = post.AuthorName
	value["community_id"] = post.CommunityID
	value["community_name"] = post.CommunityName
	value["status"] = post.Status
	value["create_time"] = post.CreateTime.Format("2006-01-02 15:04:05")
	value["like_count"] = post.LikeCount
	value["unlike_count"] = post.UnlikeCount

	_, err := client.HMSet(key, value).Result()
	if err != nil {
		return "", err
	}

	// 设置随机TTL
	ttl := RandTTL()
	_, err = client.Expire(key, ttl).Result()
	if err != nil {
		return "", err
	}

	return ttl.String(), err

}

func RandTTL() time.Duration {
	return time.Second * time.Duration(rand.Intn(MaxExpireTime-MinExpireTime)+MinExpireTime)
}

func DelayTTL(key string) {
	ttl, _ := client.TTL(key).Result()
	client.Expire(key, ttl+RandTTL())
}
