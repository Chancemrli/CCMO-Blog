package logic

import (
	"context"
	"encoding/json"

	"CCMO/gozero/blog/application/article/api/internal/code"
	"CCMO/gozero/blog/application/article/api/internal/svc"
	"CCMO/gozero/blog/application/article/api/internal/types"
	"CCMO/gozero/blog/application/article/rpc/pb"
	"CCMO/gozero/blog/application/egg/rpc/egg"
	"CCMO/gozero/blog/pkg/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

const minContentLen = 80

type PublishLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishLogic {
	return &PublishLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishLogic) Publish(req *types.PublishRequest) (*types.PublishResponse, error) {
	if len(req.Title) == 0 {
		return nil, code.ArtitleTitleEmpty
	}
	if len(req.Content) < minContentLen {
		return nil, code.ArticleContentTooFewWords
	}
	if len(req.Cover) == 0 {
		return nil, code.ArticleCoverEmpty
	}
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value error: %v", err)
		return nil, xcode.NoLogin
	}

	pret, err := l.svcCtx.ArticleRPC.Publish(l.ctx, &pb.PublishRequest{
		UserId:      userId,
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		Cover:       req.Cover,
	})
	if err != nil {
		logx.Errorf("l.svcCtx.ArticleRPC.Publish req: %v userId: %d error: %v", req, userId, err)
		return nil, err
	}

	_, err = l.svcCtx.EggRPC.Comment(l.ctx, &egg.EggRequest{
		Content: req.Content,
	})

	if err != nil {
		logx.Errorf("l.svcCtx.EggRPC.Comment error: %v", err)
		return nil, err
	}

	return &types.PublishResponse{ArticleId: pret.ArticleId}, nil
}
