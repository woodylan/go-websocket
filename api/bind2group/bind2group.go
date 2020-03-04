package bind2group

import (
	"encoding/json"
	"go-websocket/api"
	"go-websocket/define/retcode"
	"go-websocket/servers"
	"net/http"
)

type Controller struct {
}

type inputData struct {
	ClientId  string `json:"clientId" validate:"required"`
	GroupName string `json:"groupName" validate:"required"`
	UserId    string `json:"userId"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	var inputData inputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := api.Validate(inputData)
	if (err != nil) {
		api.Render(w, retcode.FAIL, err.Error(), []string{})
		return
	}

	if len(inputData.ClientId) > 0 && len(inputData.GroupName) > 0 {
		systemId := r.Header.Get("systemId")
		servers.AddClient2Group(systemId, inputData.GroupName, inputData.ClientId, inputData.UserId)
	}

	api.Render(w, retcode.SUCCESS, "success", []string{})
}
