package servers

import (
	"encoding/json"
	"github.com/astaxie/beego/config"
	"github.com/gorilla/websocket"
	"go-websocket/tools/readconfig"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testServer struct {
	*httptest.Server
	ClientURL string
}

type connectMessage struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data connectData `json:"data"`
}

type connectData struct {
	ClientId string `json:"clientId"`
}

func newServer(t *testing.T) *testServer {
	var s testServer

	readconfig.ConfigData, _ = config.NewConfig("ini", "../configs/config.ini")

	websocketHandler := &Controller{}
	s.Server = httptest.NewServer(http.HandlerFunc(websocketHandler.Run))
	s.ClientURL = "ws" + strings.TrimPrefix(s.Server.URL, "http")
	return &s
}

func TestConnect(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	ws, _, err := websocket.DefaultDialer.Dial(s.ClientURL, nil)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer ws.Close()

	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}
	render := connectMessage{}
	if err := json.Unmarshal(message, &render); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if render.Code != 0 {
		t.Error("response Code err")
		t.FailNow()
	}

	if render.Msg != "success" {
		t.Error("response Msg err")
		t.FailNow()
	}

	if len(render.Data.ClientId) <= 0 {
		t.Error("client id empty")
		t.FailNow()
	}
}
