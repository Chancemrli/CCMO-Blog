package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ReplyCountModel = (*customReplyCountModel)(nil)

type (
	// ReplyCountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReplyCountModel.
	ReplyCountModel interface {
		replyCountModel
		withSession(session sqlx.Session) ReplyCountModel
	}

	customReplyCountModel struct {
		*defaultReplyCountModel
	}
)

// NewReplyCountModel returns a model for the database table.
func NewReplyCountModel(conn sqlx.SqlConn) ReplyCountModel {
	return &customReplyCountModel{
		defaultReplyCountModel: newReplyCountModel(conn),
	}
}

func (m *customReplyCountModel) withSession(session sqlx.Session) ReplyCountModel {
	return NewReplyCountModel(sqlx.NewSqlConnFromSession(session))
}
