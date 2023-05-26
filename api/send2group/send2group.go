package send2group

import (
	"encoding/json"
	"net/http"

	"github.com/woodylan/go-websocket/api"
	"github.com/woodylan/go-websocket/define/retcode"
	"github.com/woodylan/go-websocket/servers"
)

type Controller struct {
}

type inputData struct {
	SendUserId string          `json:"sendUserId"`
	GroupName  string          `json:"groupName" validate:"required"`
	Code       int             `json:"code"`
	Msg        string          `json:"msg"`
	Data       json.RawMessage `json:"data"`
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
	dataStr := string(inputData.Data)
	systemId := r.Header.Get("SystemId")
	messageId := servers.SendMessage2Group(systemId, inputData.SendUserId, inputData.GroupName, inputData.Code, inputData.Msg, &dataStr)

	api.Render(w, retcode.SUCCESS, "success", map[string]string{
		"messageId": messageId,
	})
	return
}
