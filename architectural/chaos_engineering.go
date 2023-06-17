package architectural

/*
“If something hurts, do it more often!”

Chaos Engineering
	如果问题经常发生，人们就会学习和思考解决它的方法
	Chaos under control
		Terminate host
		Inject latency
		Inject failure
Chaos Engineering 原则
	Build a Hypothesis around Steady State Behavior
	Vary Real-world Events
	Run Experiments in Production
	Automate Experiments to Run Continuously
	Minimize Blast Radius

	http://principlesofchaos.org
相关开源项⽬目
	https://github.com/Netflix/chaosmonkey
蔡超老师的开源项目：https://github.com/easierway/service_decorators/blob/master/README.md
	decorators 模式

	强烈推荐：使用Go开发分布式服务做Microservice
	功能：
	利用声明的方式，把核心逻辑包起来：即把编码变成仅声明，自己的service就具有了这些功能
	利用 ChaosEngineeringDecorator，往自己的service注入超时/设定的错误（runtime时）
	...
*/
