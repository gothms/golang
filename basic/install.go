package basic

/*
新建
	GOROOT
	GOPATH
Path
	%GOROOT%\bin
	%GOPATH%\bin
	%USERPROFILE%\go\bin
env
	GO111MODULE=auto：$ go env -w GO111MODULE=on
		$ go env -u GO111MODULE：恢复默认值
	GOMODCACHE=E:\gospace\pkg\mod
	GOPROXY=https://goproxy.cn：$ go env -w GOPROXY=https://goproxy.cn,direct
		windows 镜像设置：
		阿里：go env -w GOPROXY="https://mirrors.aliyun.com/goproxy/"
		七牛云：go env -w GOPROXY="https://goproxy.cn"
gopath：直接手动在 %GOPATH% 创建即可(不要\go)
	mkdir %GOPATH%\go
	mkdir %GOPATH%\go\bin：用于存放编译后生成的可执行文件
	mkdir %GOPATH%\go\pkg：存放编译后生成的包文件，即go install后生成的文件。并且，如果是通过go mod管理包，那么包也会下载到这个文件夹
	mkdir %GOPATH%\go\src：用于存放项目的源码。在不启用go mod的情况下只能将源码放在该文件夹中

	mkdir %USERPROFILE%\go
	mkdir %USERPROFILE%\go\bin
	mkdir %USERPROFILE%\go\pkg
	mkdir %USERPROFILE%\go\src

$GOPATH\bin\go-torch.exe 全路径示例
	E:\gothmslee\bin\go-torch.exe
参考
	https://juejin.cn/post/7044908192882524196
*/
