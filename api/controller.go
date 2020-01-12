package api

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go-websocket/define/retcode"
	"io"
	"net/http"
)

type RetData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Date interface{} `json:"data"`
}

func ConnRender(conn *websocket.Conn, date interface{}) (err error) {
	err = conn.WriteJSON(RetData{
		Code: retcode.SUCCESS,
		Msg:  "success",
		Date: date,
	})

	return
}

func Render(w http.ResponseWriter, code int, msg string, data interface{}) (str string) {
	var retData RetData

	retData.Code = code
	retData.Msg = msg
	retData.Date = data

	retJson, _ := json.Marshal(retData)
	str = string(retJson)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, str)
	return
}
