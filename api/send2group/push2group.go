package send2group

import (
	"encoding/json"
	"go-websocket/api"
	"go-websocket/define/retcode"
	"go-websocket/servers/server"
	"net/http"
)

type Controller struct {
}

type inputData struct {
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
	var inputData inputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server.SendMessage2Group(&inputData.GroupName, &inputData.Message)

	api.Render(w, retcode.SUCCESS, "success", []string{})
	return
}
