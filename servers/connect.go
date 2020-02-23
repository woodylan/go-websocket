package servers

import (
	"github.com/gorilla/websocket"
	"go-websocket/api"
	"go-websocket/define"
	"go-websocket/define/retcode"
	"go-websocket/pkg/redis"
	"go-websocket/tools/util"
	"log"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type Controller struct {
}

type renderData struct {
	ClientId string `json:"clientId"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		http.NotFound(w, r)
		return
	}

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//解析参数
	systemId := r.FormValue("systemId")
	if len(systemId) == 0 {
		_ = Render(conn, "", "", retcode.SYSTEM_ID_ERROR, "系统ID不能为空", []string{})
		_ = conn.Close()
		return
	}

	//判断系统ID是否存在
	jsonValue, err := redis.Get(define.REDIS_PREFIX_ACCOUNT_INFO + systemId)
	if err != nil || len(jsonValue) == 0 {
		_ = Render(conn, "", "", retcode.SYSTEM_ID_ERROR, "系统ID不正确", []string{})
		_ = conn.Close()
		return
	}

	clientId := util.GenClientId()

	clientSocket := NewClient(clientId, systemId, conn)

	Manager.AddClient2SystemClient(&systemId, clientSocket)

	//读取客户端消息
	clientSocket.Read()

	if err = api.ConnRender(conn, renderData{ClientId: clientId}); err != nil {
		_ = conn.Close()
		return
	}

	// 用户连接事件
	Manager.Connect <- clientSocket
}
