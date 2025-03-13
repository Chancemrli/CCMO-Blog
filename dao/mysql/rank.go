package mysql

import (
	"bluebell_backend/models"
	"fmt"
)

func GetRankByScore(topNum int64) ([]models.PostInScoreRank, error) {
	posts := make([]models.PostInScoreRank, 0)
	sqlstr := fmt.Sprintf("select post_id, title, score from post order by score desc limit %d", topNum)
	if err := db.Select(&posts, sqlstr); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return posts, nil
}
