package controller

import (
	"bluebell_backend/logic"
	"bluebell_backend/pkg/errcode"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LikePostHandler(ctx *gin.Context) {
	// 获取用户
	id, err := getCurrentUserID(ctx)
	userID := strconv.Itoa(int(id))
	// 获取点赞对象
	postID := ctx.Param("id")
	if err != nil {
		status := errcode.GetStatus(errcode.ProgramError)
		ResponseWithStatus(ctx, status, nil)
		return
	}

	// 点赞
	likeCount, status := logic.PostLikeLogic(userID, postID)

	ResponseWithStatus(ctx, status, likeCount)

}

func UnlikePostHandler(ctx *gin.Context) {
	// 获取用户
	id, err := getCurrentUserID(ctx)
	userID := strconv.Itoa(int(id))
	// 获取点踩对象
	postID := ctx.Param("id")
	if err != nil {
		status := errcode.GetStatus(errcode.ProgramError)
		ResponseWithStatus(ctx, status, nil)
		return
	}

	// 点踩
	likeCount, status := logic.PostUnlikeLogic(userID, postID)

	ResponseWithStatus(ctx, status, likeCount)

}
