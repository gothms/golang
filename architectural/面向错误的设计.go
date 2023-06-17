package architectural

/*
“Once you accept that failures will happen,
you have the ability to design your system’s reaction to the failures.”

1.隔离
	隔离错误：
		设计：如 Micro kernel
		部署：如 Microservice
	重用 vs 隔离：
		逻辑结构的重用 vs 部署结构的隔离
2.冗余
	单点失效：
		限流
	慢响应：A quick rejection is better than a slow response.
		不要无休止的等待：给阻塞操作都加上一个期限
	错误传递：
		断路器
*/
