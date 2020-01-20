package readconfig

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInitConfig(t *testing.T) {
	Convey("读取配置文件", t, func() {
		err := InitConfig()
		So(err, ShouldBeNil)
	})
}
