package api

import (
	"encoding/json"
)

type RetData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Date interface{} `json:"data"`
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

