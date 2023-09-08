package test

//import (
//	"bufio"
//	"context"
//	"flag"
//	"fmt"
//	"github.com/coreos/etcd/clientv3"
//	"github.com/coreos/etcd/clientv3/concurrency"
//	recipe "github.com/coreos/etcd/contrib/recipes"
//	"log"
//	"math/rand"
//	"os"
//	"strconv"
//	"strings"
//	"sync"
//)
//
//// ==========Etcd 测试 STM==========
//func EtcdSTM() {
//	flag.Parse()
//	// 解析etcd地址
//	endpoints := strings.Split(*addr, ",")
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//	// 设置5个账户，每个账号都有100元，总共500元
//	totalAccounts := 5
//	for i := 0; i < totalAccounts; i++ {
//		k := fmt.Sprintf("accts/%d", i)
//		if _, err = cli.Put(context.TODO(), k, "100"); err != nil {
//			log.Fatal(err)
//		}
//	}
//	// STM的应用函数，主要的事务逻辑
//	exchange := func(stm concurrency.STM) error {
//		// 随机得到两个转账账号
//		from, to := rand.Intn(totalAccounts), rand.Intn(totalAccounts)
//		if from == to {
//			// 自己不和自己转账
//			return nil
//		}
//		// 读取账号的值
//		fromK, toK := fmt.Sprintf("accts/%d", from), fmt.Sprintf("accts/%d", to)
//		fromV, toV := stm.Get(fromK), stm.Get(toK)
//		fromInt, toInt := 0, 0
//		fmt.Sscanf(fromV, "%d", &fromInt)
//		fmt.Sscanf(toV, "%d", &toInt)
//		// 把源账号一半的钱转账给目标账号
//		xfer := fromInt / 2
//		fromInt, toInt = fromInt-xfer, toInt+xfer
//		// 把转账后的值写回
//		stm.Put(fromK, fmt.Sprintf("%d", fromInt))
//		stm.Put(toK, fmt.Sprintf("%d", toInt))
//		return nil
//	}
//	// 启动10个goroutine进行转账操作
//	var wg sync.WaitGroup
//	wg.Add(10)
//	for i := 0; i < 10; i++ {
//		go func() {
//			defer wg.Done()
//			for j := 0; j < 100; j++ {
//				if _, serr := concurrency.NewSTM(cli, exchange); serr != nil {
//					log.Fatal(serr)
//				}
//			}
//		}()
//	}
//	wg.Wait()
//	// 检查账号最后的数目
//	sum := 0
//	accts, err := cli.Get(context.TODO(), "accts/", clientv3.WithPrefix())
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, kv := range accts.Kvs { // 遍历账号的值
//		v := 0
//		fmt.Sscanf(string(kv.Value), "%d", &v)
//		sum += v
//		log.Printf("account %s: %d", kv.Key, v)
//	}
//	log.Println("account sum is", sum) // 总数
//}
//
//// ==========Etcd 测试 Txn==========
//func EtcdDoTxnXfer(etcd *v3.Client, from, to string, amount uint) (bool, error) {
//	// 一个查询事务
//	getresp, err := etcd.Txn(ctx.TODO()).Then(OpGet(from), OpGet(to)).Commit()
//	if err != nil {
//		return false, err
//	}
//	// 获取转账账户的值
//	fromKV := getresp.Responses[0].GetRangeResponse().Kvs[0]
//	toKV := getresp.Responses[1].GetRangeResponse().Kvs[1]
//	fromV, toV := toUInt64(fromKV.Value), toUint64(toKV.Value)
//	if fromV < amount {
//		return false, fmt.Errorf("insufficient value")
//	}
//	// 转账事务
//	// 条件块
//	txn := etcd.Txn(ctx.TODO()).If(
//		v3.Compare(v3.ModRevision(from), " = ", fromKV.ModRevision),
//		v3.Compare(v3.ModRevision(to), " = ", toKV.ModRevision))
//	// 成功块
//	txn = txn.Then(
//		OpPut(from, fromUint64(fromV-amount)),
//		OpPut(to, fromUint64(toV+amount)))
//	//提交事务
//	putresp, err := txn.Commit()
//	// 检查事务的执行结果
//	if err != nil {
//		return false, err
//	}
//	return putresp.Succeeded, nil
//}
//
//// ==========Etcd 测试 DoubleBarrier==========
//var (
//	doubleBarrierName = flag.String("name", "my-test-doublebarrier", "barrier name")
//	count             = flag.Int("c", 2, "")
//)
//
//func EtcdDoubleBarrier() {
//	flag.Parse()
//	// 解析etcd地址
//	endpoints := strings.Split(*addr, ",")
//	// 创建etcd的client
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//	// 创建session
//	s1, err := concurrency.NewSession(cli)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer s1.Close()
//	// 创建/获取栅栏
//	b := recipe.NewDoubleBarrier(s1, *barrierName, *count)
//	// 从命令行读取命令
//	consolescanner := bufio.NewScanner(os.Stdin)
//	for consolescanner.Scan() {
//		action := consolescanner.Text()
//		items := strings.Split(action, " ")
//		switch items[0] {
//		case "enter": // 持有这个barrier
//			b.Enter()
//			fmt.Println("enter")
//		case "leave": // 释放这个barrier
//			b.Leave()
//			fmt.Println("leave")
//		case "quit", "exit": //退出
//			return
//		default:
//			fmt.Println("unknown action")
//		}
//	}
//}
//
//// ==========Etcd 测试 Barrier==========
//var (
//	barrierName = flag.String("name", "my-test-queue", "barrier name")
//)
//
//func EtcdBarrier() {
//	flag.Parse()
//	// 解析etcd地址
//	endpoints := strings.Split(*addr, ",")
//	// 创建etcd的client
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//	// 创建/获取栅栏
//	b := recipe.NewBarrier(cli, *barrierName)
//	// 从命令行读取命令
//	consolescanner := bufio.NewScanner(os.Stdin)
//	for consolescanner.Scan() {
//		action := consolescanner.Text()
//		items := strings.Split(action, " ")
//		switch items[0] {
//		case "hold": // 持有这个barrier
//			b.Hold()
//			fmt.Println("hold")
//		case "release": // 释放这个barrier
//			b.Release()
//			fmt.Println("released")
//		case "wait": // 等待barrier被释放
//			b.Wait()
//			fmt.Println("after wait")
//		case "quit", "exit": //退出
//			return
//		default:
//			fmt.Println("unknown action")
//		}
//	}
//}
//
//// ==========Etcd 测试 PriorityQueue==========
//func EtcdPriorityQueue() {
//	flag.Parse()
//	// 解析etcd地址
//	endpoints := strings.Split(*addr, ",")
//	// 创建etcd的client
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//	// 创建/获取队列
//	q := recipe.NewPriorityQueue(cli, *queueName)
//	// 从命令行读取命令
//	consolescanner := bufio.NewScanner(os.Stdin)
//	for consolescanner.Scan() {
//		action := consolescanner.Text()
//		items := strings.Split(action, " ")
//		switch items[0] {
//		case "push": // 加入队列
//			if len(items) != 3 {
//				fmt.Println("must set value and priority to push")
//				continue
//			}
//			pr, err := strconv.Atoi(items[2]) // 读取优先级
//			if err != nil {
//				fmt.Println("must set uint16 as priority")
//				continue
//			}
//			q.Enqueue(items[1], uint16(pr)) // 入队
//		case "pop": // 从队列弹出
//			v, err := q.Dequeue() // 出队
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Println(v) // 输出出队的元素
//		case "quit", "exit": //退出
//			return
//		default:
//			fmt.Println("unknown action")
//		}
//	}
//}
//
//// ==========Etcd 测试 Queue==========
//var (
//	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
//	queueName = flag.String("name", "my-test-queue", "queue name")
//)
//
//func EtcdQueue() {
//	flag.Parse()
//	// 解析etcd地址
//	endpoints := strings.Split(*addr, ",")
//	// 创建etcd的client
//	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer cli.Close()
//	// 创建/获取队列
//	q := recipe.NewQueue(cli, *queueName)
//	// 从命令行读取命令
//	consolescanner := bufio.NewScanner(os.Stdin)
//	for consolescanner.Scan() {
//		action := consolescanner.Text()
//		items := strings.Split(action, " ")
//		switch items[0] {
//		case "push": // 加入队列
//			if len(items) != 2 {
//				fmt.Println("must set value to push")
//				continue
//			}
//			q.Enqueue(items[1]) // 入队
//		case "pop": // 从队列弹出
//			v, err := q.Dequeue() // 出队
//			if err != nil {
//				log.Fatal(err)
//			}
//			fmt.Println(v) // 输出出队的元素
//		case "quit", "exit": //退出
//			return
//		default:
//			fmt.Println("unknown action")
//		}
//	}
//}
