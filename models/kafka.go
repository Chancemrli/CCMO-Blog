package models

type UserTracking struct {
	UserID uint64 `json:"user_id" db:"user_id"`
	PostID uint64 `json:"post_id" db:"post_id"`
	// -1:unlike, 1:like
	Opt int64 `json:"opt" db:"opt"`
}
