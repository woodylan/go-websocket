package push2group

import (
	"encoding/json"
	"go-websocket/api"
	"go-websocket/servers/server"
	"io"
	"net/http"
)

type Controller struct {
}

type pushToGroupInputData struct {
	GroupName string `json:"groupName"`
	Message   string `json:"message"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
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

	server.SendMessage2Group(inputData.GroupName, inputData.Message)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, api.Render(0, "success", []string{}))

	return
}
