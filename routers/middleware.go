package routers

import (
	"github.com/woodylan/go-websocket/api"
	"github.com/woodylan/go-websocket/define"
	"github.com/woodylan/go-websocket/define/retcode"
	"github.com/woodylan/go-websocket/pkg/etcd"
	"github.com/woodylan/go-websocket/servers"
	"github.com/woodylan/go-websocket/tools/util"
	"net/http"
)

func AccessTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//检查header是否设置SystemId
		systemId := r.Header.Get("SystemId")
		if len(systemId) == 0 {
			api.Render(w, retcode.FAIL, "系统ID不能为空", []string{})
			return
		}

		//判断是否被注册
		if util.IsCluster() {
			resp, err := etcd.Get(define.ETCD_PREFIX_ACCOUNT_INFO + systemId)
			if err != nil {
				api.Render(w, retcode.FAIL, "etcd服务器错误", []string{})
				return
			}

			if resp.Count == 0 {
				api.Render(w, retcode.FAIL, "系统ID无效", []string{})
				return
			}
		} else {
			if _, ok := servers.SystemMap.Load(systemId); !ok {
				api.Render(w, retcode.FAIL, "系统ID无效", []string{})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
