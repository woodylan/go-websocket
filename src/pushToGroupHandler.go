package src

import (
	"encoding/json"
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

	b.binder.SendMessage2Group(inputData.GroupName, inputData.Message)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, render(0, "success", []string{}))

	return
}
