package practice

import (
	"fmt"
	"log"
	"sync"
	"time"
)

/*
条件变量sync.Cond （上）

条件变量（conditional variable）
	条件变量是基于互斥锁的，它必须有互斥锁的支撑才能发挥作用
	条件变量并不是被用来保护临界区和共享资源的，它是用于协调想要访问共享资源的那些线程的
	当共享资源的状态发生变化时，它可以被用来通知被互斥锁阻塞的线程
示例
	两个人在共同执行一项秘密任务，这需要在不直接联系和见面的前提下进行。我需要向一个信箱里放置情报，你需要从这个信箱中获取情报
	这个信箱就相当于一个共享资源，而我们就分别是进行写操作的线程和进行读操作的线程
		如果我在放置的时候发现信箱里还有未被取走的情报，那就不再放置，而先返回
		另一方面，如果你在获取的时候发现信箱里没有情报，那也只能先回去了
		这就相当于写的线程或读的线程阻塞的情况
	虽然我们俩都有信箱的钥匙，但是同一时刻只能有一个人插入钥匙并打开信箱，这就是锁的作用了
		更何况咱们俩是不能直接见面的，所以这个信箱本身就可以被视为一个临界区
	我们又想了一个计策，各自雇佣了一个不起眼的小孩儿
		如果早上七点有一个戴红色帽子的小孩儿从你家楼下路过，那么就意味着信箱里有了新情报
		另一边，如果上午九点有一个戴蓝色帽子的小孩儿从我家楼下路过，那就说明你已经从信箱中取走了情报
		这两个戴不同颜色帽子的小孩儿就相当于条件变量，在共享资源的状态产生变化的时候，起到了通知的作用
	条件变量在这里的最大优势就是在效率方面的提升
		当共享资源的状态不满足条件的时候，想操作它的线程再也不用循环往复地做检查了，只要等待通知就好了

问题：条件变量怎样与互斥锁配合使用？
典型回答：条件变量的初始化离不开互斥锁，并且它的方法有的也是基于互斥锁的
	条件变量提供的方法有三个：等待通知（wait）、单发通知（signal）和广播通知（broadcast）
	在利用条件变量等待通知的时候，需要在它基于的那个互斥锁保护下进行
	而在进行单发通知或广播通知的时候，却是恰恰相反的，需要在对应的互斥锁解锁之后再做这两种操作
问题解析
	示例
		var mailbox uint8	// 代表信箱，是uint8类型的。值为0则表示信箱中没有情报，值为1时则说明信箱中有情报
		var lock sync.RWMutex	// 被视为信箱上的那把锁
		sendCond := sync.NewCond(&lock)	// sendCond和recvCond，都是*sync.Cond类型的，同时也都是由sync.NewCond函数来初始化的
		recvCond := sync.NewCond(lock.RLocker())
	sync.Cond
		sync.Cond类型并不是开箱即用的，只能利用sync.NewCond函数创建它的指针值
		这个函数需要一个sync.Locker类型的参数值
		条件变量是基于互斥锁的，它必须有互斥锁的支撑才能够起作用。因此，这里的参数值是不可或缺的，它会参与到条件变量的方法实现当中
	sync.Locker 接口
		声明中只包含了两个方法定义，即：Lock()和Unlock()
		sync.Mutex类型和sync.RWMutex类型都拥有Lock方法和Unlock方法，只不过它们都是指针方法
		因此，这两个类型的指针类型才是sync.Locker接口的实现类型
	sendCond
		lock变量的Lock方法和Unlock方法分别用于对其中写锁的锁定和解锁，它们与sendCond变量的含义是对应的
		sendCond是专门为放置情报而准备的条件变量，向信箱里放置情报，可以被视为对共享资源的写操作
	recvCond
		recvCond变量代表的是专门为获取情报而准备的条件变量
		在这里，我们暂且把获取情报看做是对共享资源的读操作
		因此，为了初始化recvCond这个条件变量，我们需要的是lock变量中的读锁，并且还需要是sync.Locker类型的
		可是，lock变量中用于对读锁进行锁定和解锁的方法却是RLock和RUnlock，它们与sync.Locker接口中定义的方法并不匹配
		lock.RLocker()的返回值，拥有Lock方法和Unlock方法，它们分别调用lock变量的RLock方法和RUnlock方法（见源码）
	四个变量
		一个是代表信箱的mailbox，一个是代表信箱上的锁的lock
		代表了蓝帽子小孩儿的sendCond，以及代表了红帽子小孩儿的recvCond
	示例：SyncCond()
		只要条件不满足，我就会通过调用sendCond变量的Wait方法，去等待你的通知，只有在收到通知之后我才会再次检查信箱
		当我需要通知你的时候，我会调用recvCond变量的Signal方法
	条件变量的基本使用规则
		利用条件变量可以实现单向的通知，而双向的通知则需要两个条件变量

总结
	条件变量是基于互斥锁的一种同步工具，它必须有互斥锁的支撑才能发挥作用
	条件变量可以协调那些想要访问共享资源的线程
	当共享资源的状态发生变化时，它可以被用来通知被互斥锁阻塞的线程

思考
	*sync.Cond类型的值可以被传递吗？那sync.Cond类型的值呢？
*/

func SyncCond() {
	var mailbox uint8     // 代表信箱，是uint8类型的。值为0则表示信箱中没有情报，值为1时则说明信箱中有情报
	var lock sync.RWMutex // lock 代表信箱上的锁
	// sendCond 代表专用于发信的条件变量
	sendCond := sync.NewCond(&lock) // sendCond和recvCond，都是*sync.Cond类型的，同时也都是由sync.NewCond函数来初始化的
	// recvCond 代表专用于收信的条件变量
	recvCond := sync.NewCond(lock.RLocker())
	//rl := lock.RLocker()
	//(*sync.RWMutex)(rl).RLock() // 为什么源码能转？请看下面示例
	//fmt.Println(&lock == rw)
	fmt.Printf("%T, %T\n", &lock, lock.RLocker())

	// sign 用于传递演示完成的信号。
	sign := make(chan struct{}, 3)
	max := 5
	go func(max int) { // 用于发信。
		defer func() {
			sign <- struct{}{}
		}()
		for i := 1; i <= max; i++ {
			time.Sleep(time.Millisecond * 500)
			lock.Lock()        // 持有信箱上的锁，并且有打开信箱的权利，而不是锁上这个锁
			for mailbox == 1 { // 检查mailbox变量的值是否等于1，也就是说，要看看信箱里是不是还存有情报
				sendCond.Wait() // 如果还有情报，那么我就回家去等蓝帽子小孩儿了
			}
			log.Printf("sender [%d]: the mailbox is empty.", i)
			mailbox = 1 // 如果信箱里没有情报，那么我就把新情报放进去
			log.Printf("sender [%d]: the letter has been sent.", i)
			lock.Unlock()     // 关上信箱、锁上锁，然后离开
			recvCond.Signal() // 及时地通知红帽子小孩儿“信箱里已经有新情报了”
		}
	}(max)
	go func(max int) { // 用于收信。
		defer func() {
			sign <- struct{}{}
		}()
		for j := 1; j <= max; j++ {
			time.Sleep(time.Millisecond * 500)
			lock.RLock()
			for mailbox == 0 {
				recvCond.Wait()
			}
			log.Printf("receiver [%d]: the mailbox is full.", j)
			mailbox = 0
			log.Printf("receiver [%d]: the letter has been received.", j)
			lock.RUnlock()
			sendCond.Signal()
		}
	}(max)

	<-sign
	<-sign
}

// Locker_ 示例对外暴漏接口，且接口实体的权限控制不对外
type Locker_ interface {
	Lock_()
}
type RWMutex_ struct{}

func (r *RWMutex_) Lock_() {}
func (r *RWMutex_) RLocker_() Locker_ { // 最精彩的函数
	return (*rlocker_)(r) // 转换为 *rlocker_ 类型，并返回接口，只对外暴漏接口
}

type rlocker_ RWMutex_

func (r *rlocker_) Lock_() {
	(*RWMutex_)(r).Lock_() // 可以类型转换，因为是 rlocker_ 类型
}
func RWMutexTypeTest() {
	//rw := RWMutex_{}
	//rl := rw.RLocker_()
	//(*RWMutex_)(rl).Lock_() // 不可以类型转换，因为是 Locker_ 接口类型
}
