package architectural

/*
“A priori prediction of all failure modes is not possible.”

健康检查
	注意僵尸进程
		池化资源耗尽
		死锁
	Let it Crash!
		recover()
		重启
构建可恢复的系统
	拒绝单体系统
	面向错误和恢复的设计
		在依赖服务不可用时，可以继续存活
		快速启动
		无状态
与客户端协商
	服务器：“我太忙了，请慢点发送数据”
	Client：“好，我一分钟后再发送”
*/
