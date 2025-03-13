package controller

import (
	"bluebell_backend/logic"
	"bluebell_backend/pkg/errcode"

	"github.com/gin-gonic/gin"
)

func RankSystemHandler(ctx *gin.Context) {
	rank, status := logic.RankByLikeCountLogic()

	ctx.JSON(int(errcode.GetHttpCode(status.Code)), gin.H{
		"status": status,
		"data":   rank,
	})
}
