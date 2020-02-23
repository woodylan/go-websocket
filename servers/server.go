package servers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go-websocket/define"
	"go-websocket/pkg/redis"
	"go-websocket/tools/util"
	"log"
	"net/http"
	"time"
)

//channel通道
var ToClientChan chan clientInfo

//channel通道结构体
type clientInfo struct {
	ClientId   *string
	SendUserId *string
	MessageId  *string
	Code       int
	Msg        string
	Data       *interface{}
}

type RetData struct {
	MessageId  string      `json:"messageId"`
	SendUserId string      `json:"sendUserId"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

// 心跳间隔
var heartbeatInterval = 25 * time.Second

func init() {
	ToClientChan = make(chan clientInfo, 1000)
}

var Manager = NewClientManager() // 管理者

func StartWebSocket() {
	websocketHandler := &Controller{}
	http.HandleFunc("/ws", websocketHandler.Run)

	go Manager.Start()
}

//发送信息到指定客户端
func SendMessage2Client(clientId *string, sendUserId *string, code int, msg string, data *interface{}) (messageId string) {
	messageId = util.GenUUID()
	if util.IsCluster() {
		addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(*clientId)
		if err != nil {
			_ = fmt.Errorf("%s", err)
			return
		}

		//如果是本机则发送到本机
		if isLocal {
			SendMessage2LocalClient(&messageId, clientId, sendUserId, code, msg, data)
		} else {
			//发送到指定机器
			SendRpc2Client(addr, &messageId, sendUserId, clientId, msg, data)
		}
	} else {
		//如果是单机服务，则只发送到本机
		SendMessage2LocalClient(&messageId, clientId, sendUserId, code, msg, data)
	}

	return
}

//添加客户端到分组
func AddClient2Group(systemId *string, groupName *string, clientId string) {
	//如果是集群则用redis共享数据
	if util.IsCluster() {
		//判断key是否存在
		addr, _, _, isLocal, err := util.GetAddrInfoAndIsLocal(clientId)
		if err != nil {
			_ = fmt.Errorf("%s", err)
			return
		}

		if isLocal {
			if client, err := Manager.GetByClientId(clientId); err == nil {
				//添加到本地
				Manager.AddClient2LocalGroup(util.GenGroupKey(*systemId, *groupName), client)
			} else {
				fmt.Println(err)
			}
		} else {
			//发送到指定的机器
			SendRpcBindGroup(&addr, systemId, groupName, &clientId)
		}
	} else {
		if client, err := Manager.GetByClientId(clientId); err == nil {
			//如果是单机，就直接添加到本地group了
			Manager.AddClient2LocalGroup(util.GenGroupKey(*systemId, *groupName), client)
		};
	}
}

//发送信息到指定分组
func SendMessage2Group(systemId, sendUserId, groupName *string, code int, msg string, data *interface{}) (messageId string) {
	messageId = util.GenUUID()
	if util.IsCluster() {
		//发送分组消息给指定广播
		SendGroupBroadcast(systemId, &messageId, sendUserId, groupName, code, msg, data)
	} else {
		//如果是单机服务，则只发送到本机
		Manager.SendMessage2LocalGroup(systemId, &messageId, sendUserId, groupName, code, msg, data)
	}
	return
}

//发送信息到指定系统
func SendMessage2System(systemId, sendUserId *string, code int, msg string, data interface{}) {
	messageId := util.GenUUID()
	if util.IsCluster() {
		//发送到系统广播
		SendSystemBroadcast(systemId, &messageId, sendUserId, code, msg, &data)
	} else {
		//如果是单机服务，则只发送到本机
		Manager.SendMessage2LocalSystem(systemId, &messageId, sendUserId, code, msg, &data)
	}
}

//获取分组列表
func GetOnlineList(systemId *string, groupName *string) map[string]interface{} {
	var clientList []string
	if util.IsCluster() {
		//发送到系统广播
		clientList = GetOnlineListBroadcast(systemId, groupName)
	} else {
		//如果是单机服务，则只发送到本机
		retList := Manager.GetGroupClientList(util.GenGroupKey(*systemId, *groupName))
		for _, clientItem := range retList {
			clientList = append(clientList, clientItem.ClientId)
		}
	}

	return map[string]interface{}{
		"count": len(clientList),
		"list":  clientList,
	}
}

//通过本服务器发送信息
func SendMessage2LocalClient(messageId, clientId *string, sendUserId *string, code int, msg string, data *interface{}) {
	ToClientChan <- clientInfo{ClientId: clientId, MessageId: messageId, SendUserId: sendUserId, Code: code, Msg: msg, Data: data}
	return
}

//监听并发送给客户端信息
func WriteMessage() {
	for {
		select {
		case clientInfo := <-ToClientChan:
			fmt.Println("发送到本机客户端：" + *clientInfo.ClientId)
			if conn, err := Manager.GetByClientId(*clientInfo.ClientId); err == nil && conn != nil {
				if err := Render(conn.Socket, *clientInfo.MessageId, *clientInfo.SendUserId, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
					_ = conn.Socket.Close()
					log.Println(err)
					return
				} else {
					//延长key过期时间
					_, err := redis.SetSurvivalTime(define.REDIS_CLIENT_ID_PREFIX+*clientInfo.ClientId, define.REDIS_KEY_SURVIVAL_SECONDS)
					if (err != nil) {
						log.Println(err)
					}
				}
			}
		}
	}
}

func Render(conn *websocket.Conn, messageId string, sendUserId string, code int, message string, data interface{}) error {
	return conn.WriteJSON(RetData{
		Code:       code,
		MessageId:  messageId,
		SendUserId: sendUserId,
		Msg:        message,
		Data:       data,
	})
}

//启动定时器进行心跳检测
func PingTimer() {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//发送心跳
				for clientId, conn := range Manager.AllClient() {
					if err := conn.Socket.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
						_ = conn.Socket.Close()
						Manager.DelClient(conn)
						log.Printf("发送心跳失败: %s 总连接数：%d", clientId, Manager.Count())
					}
				}
			}
		}
	}()
}
