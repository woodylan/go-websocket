package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PushToGroupHandler struct {
	binder *binder
}

type pushToGroupInputData struct {
	GroupName string `json:"groupName"`
	Message   string `json:"message"`
}

func (b *PushToGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//解析参数
	_ = r.ParseForm()
	var inputData pushToGroupInputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(inputData.GroupName) > 0 {
		if clientList, ok := b.binder.groupClientIds[inputData.GroupName]; ok {
			for _, client := range clientList {
				fmt.Println("发送消息")
				fmt.Println(client)
				//发送信息
				toClientChan <- [2]string{client, inputData.Message}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, render(0, "success", []string{}))

	return
}
