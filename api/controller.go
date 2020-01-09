package api

import (
	"encoding/json"
	"go-websocket/clientvar"
	"go-websocket/servers"
	"log"
)

type RetData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Date interface{} `json:"data"`
}

func WriteMessage() {
	for {
		select {
		case clientInfo := <-servers.ToClientChan:
			toConn, ok := clientvar.IsAlive(clientInfo[0]);
			if ok {
				err := toConn.WriteJSON(clientInfo[1]);
				if err != nil {
					go clientvar.DelClient(clientInfo[0])
					log.Println(err)
				} else {
					//todo 给redis续命
				}
			} else {
				go clientvar.DelClient(clientInfo[0])
			}
		}
	}
}

func Render(code int, msg string, data interface{}) (str string) {
	var retData RetData

	retData.Code = code
	retData.Msg = msg
	retData.Date = data

	retJson, _ := json.Marshal(retData)
	str = string(retJson)
	return
}

