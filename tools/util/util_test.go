package util

import (
	. "github.com/smartystreets/goconvey/convey"
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
