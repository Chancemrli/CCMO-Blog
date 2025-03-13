package mysql

import (
	"bluebell_backend/models"
	"bluebell_backend/pkg/util"
)

func IncrLikeCount(postID uint64) error {
	// 事务操作
	tx, err := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	postScore := models.PostScore{
		PostID: postID,
	}
	selectSQL := "select like_count, unlike_count, create_time from post where post_id = ?"
	if err = tx.QueryRow(selectSQL, postID).Scan(&postScore.LikeCount, &postScore.UnlikeCount, &postScore.CreateTime); err != nil {
		return ErrorQueryFailed
	}
	postScore.LikeCount++
	postScore.Score = int64(util.Hot(int(postScore.LikeCount), int(postScore.UnlikeCount), postScore.CreateTime))

	sqlStr := "update post set like_count = ?, score = ? where post_id = ?"

	if _, err = tx.Exec(sqlStr, postScore.LikeCount, postScore.Score, postScore.PostID); err != nil {
		return ErrorUpdateFailed
	}
	return nil
}

func IncrUnlikeCount(postID uint64) error {
	// 事务操作
	tx, err := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	postScore := models.PostScore{
		PostID: postID,
	}
	selectSQL := "select like_count, unlike_count, create_time from post where post_id = ?"
	if err = tx.QueryRow(selectSQL, postID).Scan(&postScore.LikeCount, &postScore.UnlikeCount, &postScore.CreateTime); err != nil {
		return ErrorQueryFailed
	}
	postScore.UnlikeCount++
	postScore.Score = int64(util.Hot(int(postScore.LikeCount), int(postScore.UnlikeCount), postScore.CreateTime))

	sqlStr := "update post set unlike_count = ?, score = ? where post_id = ?"

	if _, err = tx.Exec(sqlStr, postScore.UnlikeCount, postScore.Score, postScore.PostID); err != nil {
		return ErrorUpdateFailed
	}
	return nil
}

func SyncData(data models.UserTracking) error {
	var err error
	if data.Opt == 1 {
		err = IncrLikeCount(data.PostID)
	} else {
		err = IncrUnlikeCount(data.PostID)
	}

	if err != nil {
		return err
	}

	return nil
}
