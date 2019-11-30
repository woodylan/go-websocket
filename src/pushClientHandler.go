package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PushToClientHandler struct {
	binder *binder
}

type pushToClientInputData struct {
	ClientId string `json:"clientId"`
	Message  string `json:"message"`
}

func (ph *PushToClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//解析参数
	_ = r.ParseForm()
	var inputData pushToClientInputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//发送信息
	toClientChan <- [2]string{inputData.ClientId, inputData.Message}

	fmt.Println(inputData.ClientId, inputData.Message)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, render(0, "success", []string{}))

	return
}
