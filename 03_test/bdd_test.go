package test_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

/*
BDD
	Behavior Driven Development：观看视频介绍

BDD in Go
	项目：https://github.com/smartystreets/goconvey
	安装：$ go get -u github.com/smartystreets/goconvey/convey
		$ go install github.com/smartystreets/goconvey
	启动 Web UI：$GOPATH/bin/goconvey

. "github.com/smartystreets/goconvey/convey"
	. 可以不写 包名

Web 界面
	$ ~/go/bin/goconvey
	window：$ goconvey
*/

func TestSpec(t *testing.T) {
	Convey("Given 2 even numbers", t, func() {
		a, b := 2, 3
		Convey("When add the two numbers", func() {
			c := a + b
			Convey("Then the result is still even", func() {
				So(c&1, ShouldEqual, 0)
			})
		})
	})
}
