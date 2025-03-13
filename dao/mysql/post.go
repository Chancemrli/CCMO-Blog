package mysql

import (
	"bluebell_backend/models"
	"bluebell_backend/pkg/util"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.PostID, post.Title,
		post.Content, post.AuthorId, post.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// GetPostByID
func GetPostByID(idStr string) (post *models.ApiPostDetail, err error) {
	post = new(models.ApiPostDetail)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time, like_count, unlike_count
	from post
	where post_id = ?`
	err = db.Get(post, sqlStr, idStr)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)`
	// 动态填充id
	query, args, err := sqlx.In(sqlStr, ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

func GetPostList() (posts []*models.ApiPostDetail, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	limit 2
	`
	posts = make([]*models.ApiPostDetail, 0, 2)
	err = db.Select(&posts, sqlStr)
	return

}

func UpdatePost(post models.UpdatePost) error {
	// 使用预编译语句以避免 SQL 注入
	sqlStr := "update post set title = ?, content = ?, status = ? where post_id = ?"

	// 执行数据库更新
	if _, err := db.Exec(sqlStr, post.Title, post.Content, post.Status, post.PostID); err != nil {
		return ErrorUpdateFailed
	}

	return nil
}

type ScoreModel struct {
	PostID      string    `db:"post_id"`
	CreateTime  time.Time `db:"create_time"`
	LikeCount   int       `db:"like_count"`
	UnlikeCount int       `db:"unlike_count"`
}

func UpdateScore() {
	scoreModel := make([]ScoreModel, 0)
	sqlStr := "select post_id,create_time,like_count,unlike_count from post"
	db.Select(&scoreModel, sqlStr)

	for _, score := range scoreModel {
		totalScore := int64(util.Hot(score.LikeCount, score.UnlikeCount, score.CreateTime))
		updateSqlStr := fmt.Sprintf("update post set score = %d where post_id = '%s'", totalScore, score.PostID)
		db.Exec(updateSqlStr)
	}
}
