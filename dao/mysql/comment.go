package mysql

import (
	"bluebell_backend/models"
	"strconv"

	"go.uber.org/zap"
)

func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `insert into comment(
	comment_id, content, post_id, author_id, parent_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, comment.CommentID, comment.Content, comment.PostID,
		comment.AuthorID, comment.ParentID)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetCommentListByIDs(postID string) (commentList []*models.Comment, err error) {
	sqlStr := `select comment_id, content, post_id, author_id, parent_id, create_time
	from comment
	where post_id = ?`
	// 动态填充id
	id, _ := strconv.ParseInt(postID, 10, 64)
	err = db.Select(&commentList, sqlStr, id)
	if err != nil {
		return
	}

	// query = db.Rebind(query)
	// err = db.Select(&commentList, query, args...)
	return
}
