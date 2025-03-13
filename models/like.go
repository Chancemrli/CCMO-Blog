package models

import "time"

type PostScore struct {
	PostID      uint64    `json:"post_id" db:"post_id"`
	LikeCount   int64     `json:"like_count" db:"like_count"`
	UnlikeCount int64     `json:"unlike_count" db:"unlike_count"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	Score       int64     `json:"score" db:"score"`
}
