package logic

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/dao/redis"
	"bluebell_backend/models"
	"bluebell_backend/pkg/errcode"
	"bluebell_backend/pkg/snowflake"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var gsf singleflight.Group

func CreatePost(post *models.Post) (err error) {
	// 生成帖子ID
	postID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	post.PostID = postID
	// 创建帖子
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	_, err = mysql.GetCommunityNameByID(fmt.Sprint(post.CommunityID))
	if err != nil {
		zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
		return err
	}
	return

}

func GetPost(postID string) (post *models.ApiPostDetail, err error) {
	var GetDataFromDB = func() (interface{}, error) {
		// 从数据库获取
		post, err := mysql.GetPostByID(postID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(postID) failed", zap.String("post_id", postID), zap.Error(err))
			return nil, err
		}
		user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
			return nil, err
		}
		community, err := mysql.GetCommunityByID(fmt.Sprint(post.CommunityID))
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			return nil, err
		}
		post.AuthorName = user.UserName
		post.CommunityName = community.CommunityName

		// 设置缓存
		_, err = redis.SetHotArticle(*post)
		if err != nil {
			zap.L().Error("redis.SetHotArticle(*post) failed", zap.Error(err))
		}
		return post, nil
	}

	// 查找缓存
	post, err = redis.GetArticleByID(postID)
	if err != nil {
		return nil, err
	}
	// 缓存层没有缓存该文章
	if post.Title == "" {
		// 合并请求
		v, err, _ := gsf.Do("GetDetail:"+postID, GetDataFromDB)
		if err != nil && err != mysql.ErrorInvalidID {
			return nil, err
		}

		// 数据库为空，设置空缓存
		if err == mysql.ErrorInvalidID {
			redis.SetEmptyArticle(postID)
			return nil, err
		}
		post = v.(*models.ApiPostDetail)
	}
	// 延长缓存TTL
	redis.DelayTTL(redis.KeyHotArticlePrefix + postID)

	return post, nil
}

func GetPostList2() (data []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList()
	if err != nil {
		fmt.Println(err)
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(postList))
	for _, post := range postList {
		user, err := mysql.GetUserByID(fmt.Sprint(post.AuthorId))
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.String("author_id", fmt.Sprint(post.AuthorId)), zap.Error(err))
			continue
		}
		post.AuthorName = user.UserName
		community, err := mysql.GetCommunityByID(fmt.Sprint(post.CommunityID))
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed", zap.String("community_id", fmt.Sprint(post.CommunityID)), zap.Error(err))
			continue
		}
		post.CommunityName = community.CommunityName
		data = append(data, post)
	}
	return
}

func UpdatePostLogic(post models.UpdatePost) errcode.Status {
	// 缓存淘汰策略：写后删
	// 更新文章
	if err := mysql.UpdatePost(post); err != nil {
		return errcode.GetStatus(errcode.SQLRunTimeError)
	}
	// 删除缓存
	postID := strconv.Itoa(int(post.PostID))
	if err := redis.DelPostByID(postID); err != nil {
		return errcode.GetStatus(errcode.RedisRunTimeError)
	}
	return errcode.GetStatus(errcode.AllOK)
}
