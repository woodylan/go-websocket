package push2client

import (
	"encoding/json"
	"go-websocket/api"
	"go-websocket/servers"
	"io"
	"net/http"
)

type Push2ClientHandler struct {
}

type pushToClientInputData struct {
	ClientId string `json:"clientId"`
	Message  string `json:"message"`
}

func (ph *Push2ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	servers.SendMessage2Client(inputData.ClientId, inputData.Message)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, api.Render(0, "success", []string{}))

	return
}
