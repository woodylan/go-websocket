package send2clients

import (
	"encoding/json"
	"github.com/woodylan/go-websocket/api"
	"github.com/woodylan/go-websocket/define/retcode"
	"github.com/woodylan/go-websocket/servers"
	"net/http"
)

type Controller struct {
}

type inputData struct {
	ClientIds  []string `json:"clientIds" validate:"required"`
	SendUserId string   `json:"sendUserId"`
	Code       int      `json:"code"`
	Msg        string   `json:"msg"`
	Data       string   `json:"data"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	var inputData inputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := api.Validate(inputData)
	if err != nil {
		api.Render(w, retcode.FAIL, err.Error(), []string{})
		return
	}

	for _, clientId := range inputData.ClientIds {
		//发送信息
		_ = servers.SendMessage2Client(clientId, inputData.SendUserId, inputData.Code, inputData.Msg, &inputData.Data)
	}

	api.Render(w, retcode.SUCCESS, "success", []string{})
	return
}
