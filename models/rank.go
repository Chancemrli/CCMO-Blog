package models

type PostInScoreRank struct {
	PostID uint64 `json:"post_id" db:"post_id"`
	Title  string `json:"title" db:"title"`
	Score  int64  `json:"score" db:"score"`
}
