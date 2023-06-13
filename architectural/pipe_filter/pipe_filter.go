package pipe_filter

/*
Pipe-Filter 架构

	1.非常适合与数据处理及数据分析系统
		数据处理
		数据分析
	2.Filter 封装数据处理的功能
	3.松耦合：Filter 只跟数据（格式）耦合
	4.Pipe 用于连接 Filter 传递数据或者在异步处理过程中缓冲数据流
		进程内同步调用时，Pipe 演变为数据在方法调用间传递

Filter 和组合模式
*/

type Request interface{}
type Response interface{}
type Filter interface {
	Process(Request) (Response, error)
}
