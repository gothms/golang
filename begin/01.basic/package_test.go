package _1_basic

import (
	"fmt"
	cm "github.com/easierway/concurrent_map"
	"golang/utils"
	"testing"
)

/*
package
	1.基本复用模块单元
		以首字母大写来表明可被包外代码访问
	2.代码的package可以和所在的目录不一致
	3.同一目录里的Go文件的package要保持一致

1.包的引用
	项目的目录（或上级某个目录）加到 GOPATH
	import "项目名/包名"：$GOPATH$/项目名/包名...

2.init方法
	在main被执行前，所有依赖的package的 init 方法都会被执行
	不同包的 init 函数按照包导入的依赖关系决定执行顺序
	每个包可以有多个 init 函数
	包的每个源文件也可以有多个 init 函数，这点比较特殊

3.远程包的引用
	引入：go get -u github.com/easierway/concurrent_map
		-u：强制从网络更新远程依赖
	使用：
		导入&别名：import cm "github.com/easierway/concurrent_map"
	HTTPS：https://github.com/easierway/concurrent_map.git
		注意代码在 GitHub 上的组织形式，以适应 go get
		直接以代码路径开始，不要有 src（即不要提交 src 目录）

4.Go未解决的依赖问题
	同一环境下，不同项目使用同一包的不同版本
	无法管理对包的特定版本的依赖

5.vendor路径
	Go 1.5 release 版本中，vendor目录被添加到除了 GOPATH 和 GOROOT 之外的依赖目录查找的解决方案
	在1.6之前，需要手动设置环境变量

	查找依赖包路径的解决方案如下：
	1.当前包下的 vendor 目录
	2.向上级目录查找，直到找到 src 下的 vendor 目录
	3.在 GOPATH 下面查找依赖包
	4.在 GOROOT 目录下查找

6.常用的依赖管理工具
	godep：https://github.com/tools/godep
	glide：https://github.com/Masterminds/glide
	dep：https://github.com/golang/dep

7.glide 的使用
	参考：https://cloud.tencent.com/developer/article/1683153

	1.安装glide：go get github.com/Masterminds/glide
		报错：
			github.com/codegangsta/cli: github.com/codegangsta/cli@v1.22.13: parsing go.mod:
				   module declares its path as: github.com/urfave/cli
						   but was required as: github.com/codegangsta/cli
		解决：go mod 中添加
			//replace github.com/urfave/cli => github.com/codegangsta/cli v1.22.13	// 错误
			replace github.com/codegangsta/cli => github.com/urfave/cli v1.22.13
		安装：
			go get github.com/Masterminds/glide
			go install github.com/Masterminds/glide
	2.使用：
		初始化：
			glide
			glide init
			Y Y N S N Y
			Y Y Y P Y Y
		使用：
			glide install：生成了 vendor 文件夹
		glide install 后依赖包飘红：
			$ go env -w GO111MODULE=on
			go mod tidy
			go mod vendor：保存到vendor
			go build -mod=vendor：编译时
	3.镜像
		glide mirror 设置：https://www.jianshu.com/p/f5f42e7915ac
		修改 glide.yaml 加入：- package: golang.org/x/crypto
			$ cat glide.yaml
		$ glide mirror set golang.org/x/crypto github.com/golang/crypto
	4.GO15VENDOREXPERIMENT
		golang 1.5引入, 默认是关闭的, 通过手动设置环境变量：GO15VENDOREXPERIMENT=1开启
		golang 1.6默认开启
		goalng 1.7 vendor作为功能支持,取消GO15VENDOREXPERIMENT环境变量
	5.版本号指定规则
		https://zhuanlan.zhihu.com/p/27994151

一个完整的 glide.yaml
	package: foor
	homepage: https://github.com/qiangmzsx
	license: MIT
	owners:
	- name: qiangmzsx
	  email: qiangmzsx@hotmail.com
	  homepage: https://github.com/qiangmzsx
	# 去除包
	ignore:
	- appengine
	- golang.org/x/net
	# 排除目录
	excludeDirs:
	- node_modules
	# 导入包
	import:
	- package: github.com/astaxie/beego
	  version: 1.8.0
	- package: github.com/coocood/freecache
	- package: github.com/garyburd/redigo/redis
	- package: github.com/go-sql-driver/mysql
	- package: github.com/bitly/go-simplejson
	- package: git.oschina.net/qiangmzsx/beegofreecache
	testImport:
	- package: github.com/smartystreets/goconvey
	  subpackages:
	  - convey
*/

func TestConcurrentMap(t *testing.T) {
	m := cm.CreateConcurrentMap(99)
	m.Set(cm.StrKey("key"), 10)
	t.Log(m.Get(cm.StrKey("key"))) // 10 true
}

func init() {
	fmt.Println("init 01")
}
func init() {
	fmt.Println("init 02")
}

func TestPackage(t *testing.T) {
	sum := utils.Add(2, 3)
	t.Log(sum)
}
