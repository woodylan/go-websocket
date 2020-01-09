package bindgroup

import (
	"encoding/json"
	"fmt"
	"go-websocket/api"
	"go-websocket/servers"
	"io"
	"net/http"
)

type Controller struct {
}

type bindToGroupInputData struct {
	ClientId  string `json:"clientId"`
	GroupName string `json:"groupName"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//解析参数
	_ = r.ParseForm()
	var inputData bindToGroupInputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(inputData.ClientId) > 0 && len(inputData.GroupName) > 0 {
		servers.AddClient2Group(inputData.GroupName, inputData.ClientId)
	} else {
		fmt.Println("参数错误")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, api.Render(0, "success", []string{}))
}
