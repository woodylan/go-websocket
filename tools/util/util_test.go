package util

import (
	"github.com/astaxie/beego/config"
	. "github.com/smartystreets/goconvey/convey"
	"go-websocket/tools/readconfig"
	"testing"
)

func TestGenUUID(t *testing.T) {
	Convey("生成uuid", t, func() {
		uuid := GenUUID()
		Convey("验证长度", func() {
			So(len(uuid), ShouldBeGreaterThan, 0)
		})
	})
}

func TestIsCluster(t *testing.T) {
	Convey("验证是否为集群", t, func() {
		Convey("是集群", func() {
			readconfig.ConfigData, _ = config.NewConfigData("ini", []byte{})
			_ = readconfig.ConfigData.Set("common::cluster", "1")
			So(IsCluster(), ShouldBeTrue)
		})

		Convey("不是集群", func() {
			readconfig.ConfigData, _ = config.NewConfigData("ini", []byte{})
			_ = readconfig.ConfigData.Set("common::cluster", "0")
			So(IsCluster(), ShouldBeFalse)
		})
	})
}

func TestGetIntranetIp(t *testing.T) {
	Convey("验证IP地址", t, func() {
		ip := GetIntranetIp()
		So(len(ip), ShouldBeGreaterThan, 0)
	})
}
