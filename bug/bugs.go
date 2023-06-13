package golang

/*
go mod 飘红：
	File->Settings...->Go->Go Modules->勾选Enable Go modules integration
		Environment(也可填入 GOPROXY=https://goproxy.io)
$GOPATH/go.mod exists but should not：报错
	手动生成 go.mod，并手动添加版本号，如 go 1.20
easayjson：easy_json.go
	描述：不能使用命令
		$ goconvey
		$ easayjson -all <file>.go
	报错：goconvey : 无法将“goconvey”项识别为 cmdlet、函数、脚本文件或可运行程序的名称。请检查名称的拼写，如果包括路径，请确保路径正确，然后再试一次。
	修复：windows下运行*.ps1脚本（powershell的脚本）的时候，需要设置执行权限
		$ set-executionpolicy remotesigned
		报错：
			set-executionpolicy : 对注册表项“HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\PowerShell\1\ShellIds\Microsoft.PowerShell”的访问被拒绝。
			要更改默认(LocalMachine)作用域的执行策略，请使用“以管理员身份运行”选项启动 Windows PowerShell。
			要更改当前用户的执行策略，请运行 "Set-ExecutionPolicy -Scope CurrentUser"。
		解决：管理员身份运行 PowerShell，再执行 $ set-executionpolicy remotesigned
		$ easyjson -all json.go
		报错：
			easyjson-bootstrap234457960.go:12:3: cannot query module due to -mod=vendor
					(Go version in go.mod is at least 1.14 and vendor directory exists.)
			Bootstrap failed: exit status 1
		解决：删除 vendor
goconvey：bdd_test.go
	$ goconvey
	报错：goconvey : 无法将“goconvey”项识别为 cmdlet、函数、脚本文件或可运行程序的名称。请检查名称的拼写，如果包括路径，请确保路径正确，然后再试一次。

glide：package_test.go 3个问题
	1.glide install 报错
		[WARN]  Unable to checkout golang/utils
		[ERROR] Update failed for golang/utils: Cannot detect VCS
		[ERROR] Failed to do initial checkout of config: Cannot detect VCS

		原因：golang/utils 为本地的依赖包
	2.glide mirror 设置：https://www.jianshu.com/p/f5f42e7915ac
		glide.yaml 加入：- package: golang.org/x/crypto
		$ glide mirror set golang.org/x/crypto github.com/golang/crypto

		测试没通过
	3.glide i：报错
		[ERROR] Unable to export dependencies to vendor directory: Error moving files: exit status 1. output: �ܾ����ʡ��ƶ��ˡ�        0 ��Ŀ¼��

		未解决
		参考：https://coder55.com/article/46327
*/
