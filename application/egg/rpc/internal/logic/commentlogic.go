package logic

import (
	"context"

	"CCMO/gozero/blog/application/egg/rpc/egg"
	"CCMO/gozero/blog/application/egg/rpc/internal/model"
	"CCMO/gozero/blog/application/egg/rpc/internal/svc"
	"CCMO/gozero/blog/pkg/agent/coze/message"

	"github.com/coze-dev/coze-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type CommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentLogic {
	return &CommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentLogic) Comment(in *egg.EggRequest) (*egg.EggResponse, error) {
	// 返回评论响应
	resp, err := l.svcCtx.Agent.Workflow(map[string]interface{}{
		l.svcCtx.Config.InPut: in.Content,
	}, l.svcCtx.Config.WID)
	if err != nil {
		return nil, err
	}

	// 解析评论内容
	comment, err := message.ParseData(resp.(coze.RunWorkflowsResp))
	if err != nil {
		return nil, err
	}

	// todo:入库
	l.svcCtx.ReplyModel.Insert(l.ctx, &model.Reply{
		TargetId:      uint64(in.ArticleId),
		ReplyUserId:   l.svcCtx.Config.User.UserID,
		BeReplyUserId: uint64(in.AuthorId),
		Content:       comment,
	})

	return &egg.EggResponse{
		Comment: comment,
	}, nil
}
