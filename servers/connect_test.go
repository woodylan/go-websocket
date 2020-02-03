package servers

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
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

	if err := readconfig.InitConfig(); err != nil {
		panic(err)
	}

	websocketHandler := &Controller{}
	s.Server = httptest.NewServer(http.HandlerFunc(websocketHandler.Run))
	s.ClientURL = "ws" + strings.TrimPrefix(s.Server.URL, "http")
	return &s
}

func TestConnect(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	ws, _, err := websocket.DefaultDialer.Dial(s.ClientURL+"?systemId=publishSystem", nil)
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

	Convey("验证json解析返回的内容", t, func() {
		err := json.Unmarshal(message, &render)
		Convey("是否解析成功", func() {
			So(err, ShouldBeNil)
		})

		Convey("Code格式", func() {
			So(render.Code, ShouldEqual, 0)
		})

		Convey("Msg格式", func() {
			So(render.Msg, ShouldEqual, "success")
		})

		Convey("Client长度", func() {
			So(len(render.Data.ClientId), ShouldBeGreaterThan, 0)
		})
	})
}
