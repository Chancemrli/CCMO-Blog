package handler

import (
	"context"
	"net/http"

	"CCMO/gozero/blog/application/applet/internal/logic"
	"CCMO/gozero/blog/application/applet/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(context.Background(), "token", r.Header.Get("Authorization"))
		r = r.WithContext(ctx)
		l := logic.NewUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserInfo()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
