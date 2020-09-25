package send2client

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/woodylan/go-websocket/pkg/setting"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testServer struct {
	*httptest.Server
	ClientURL string
}

type retMessage struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func newServer(t *testing.T) *testServer {
	var s testServer
	setting.Default()

	controller := &Controller{}
	s.Server = httptest.NewServer(http.HandlerFunc(controller.Run))
	s.ClientURL = s.Server.URL + "/api/send_to_client"
	return &s
}

func TestRun(t *testing.T) {
	s := newServer(t)
	defer s.Close()

	testContent := `{"clientId":"ade447d79f6489b5","code":0,"msg":"success","data":"发送给指定客户端"}`

	resp, err := http.Post(s.ClientURL, "application/json", strings.NewReader(testContent))
	Convey("测试发送消息给指定客户端", t, func() {
		Convey("是否有报错", func() {
			So(err, ShouldBeNil)
		})
	})
	defer resp.Body.Close()

	retMessage := retMessage{}
	message, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(message, &retMessage)

	Convey("验证json解析返回的内容", t, func() {
		err := json.Unmarshal(message, &retMessage)
		Convey("是否解析成功", func() {
			So(err, ShouldBeNil)
		})

		Convey("Code格式", func() {
			So(retMessage.Code, ShouldEqual, 0)
		})

		Convey("Msg格式", func() {
			So(retMessage.Msg, ShouldEqual, "success")
		})

	})
}
