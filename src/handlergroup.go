package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BindToGroupHandler struct {
	binder *binder
}

type bindToGroupInputData struct {
	ClientId  string `json:"clientId"`
	GroupName string `json:"groupName"`
}

func (h *BindToGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		h.binder.AddClient2Group(inputData.GroupName, inputData.ClientId)
	} else {
		fmt.Println("参数错误")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, render(0, "success", []string{}))
}
