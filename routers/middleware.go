package routers

import (
	"go-websocket/api"
	"go-websocket/define"
	"go-websocket/define/retcode"
	"go-websocket/pkg/redis"
	"net/http"
)

func AccessTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//检查header是否设置token
		token := r.Header.Get("Token")
		if len(token) == 0 {
			api.Render(w, retcode.FAIL, "token未设置或无效", []string{})
			return
		}

		//校验token是否合格
		systemName, err := redis.Get(define.REDIS_PREFIX_TOKEN + token)
		if err != nil {
			api.Render(w, retcode.FAIL, "redis服务器错误", []string{})
			return
		}

		if len(systemName) == 0 {
			api.Render(w, retcode.FAIL, "token未设置或无效", []string{})
			return
		}

		next.ServeHTTP(w, r)
	})
}
