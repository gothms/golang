package test

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	recipe "github.com/coreos/etcd/contrib/recipes"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ==========Etcd 测试 RWMutex==========
func EtcdRWMutex() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	// 解析etcd地址
	endpoints := strings.Split(*addr, ",")
	// 创建etcd的client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// 创建session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := recipe.NewRWMutex(s1, *lockName)
	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		switch action {
		case "w": // 请求写锁
			testWriteLocker(m1)
		case "r": // 请求读锁
			testReadLocker(m1)
		default:
			fmt.Println("unknown action")
		}
	}
}
func testWriteLocker(m1 *recipe.RWMutex) {
	// 请求写锁
	log.Println("acquiring write lock")
	if err := m1.Lock(); err != nil {
		log.Fatal(err)
	}
	log.Println("acquired write lock")
	// 等待一段时间
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	// 释放写锁
	if err := m1.Unlock(); err != nil {
		log.Fatal(err)
	}
	log.Println("released write lock")
}
func testReadLocker(m1 *recipe.RWMutex) {
	// 请求读锁
	log.Println("acquiring read lock")
	if err := m1.RLock(); err != nil {
		log.Fatal(err)
	}
	log.Println("acquired read lock")
	// 等待一段时间
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	// 释放写锁
	if err := m1.RUnlock(); err != nil {
		log.Fatal(err)
	}
	log.Println("released read lock")
}

// ==========Etcd 测试 Mutex==========
func useMutex(cli *clientv3.Client) {
	// 为锁生成session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	m1 := concurrency.NewMutex(s1, *lockName)
	//在请求锁之前查询key
	log.Printf("before acquiring. key: %s", m1.Key())
	// 请求锁
	log.Println("acquiring lock")
	if err := m1.Lock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Printf("acquired lock. key: %s", m1.Key())
	//等待一段时间
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	// 释放锁
	if err := m1.Unlock(context.TODO()); err != nil {
		log.Fatal(err)
	}
	log.Println("released lock")
}

// ==========Etcd 测试 Locker==========
var (
	addr     = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	lockName = flag.String("name", "my-test-lock", "lock name")
)

func EtcdLocker() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	// etcd地址
	endpoints := strings.Split(*addr, ",")
	// 生成一个etcd client
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	useLock(cli) // 测试锁
}
func useLock(cli *clientv3.Client) {
	// 为锁生成session
	s1, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}
	defer s1.Close()
	//得到一个分布式锁
	locker := concurrency.NewLocker(s1, *lockName)
	// 请求锁
	log.Println("acquiring lock")
	locker.Lock()
	log.Println("acquired lock")
	// 等待一段时间
	time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
	locker.Unlock() // 释放锁
	log.Println("released lock")
}

// ==========Etcd 测试监控命令（Observe）==========
func watch(e1 *concurrency.Election, electName string) {
	ch := e1.Observe(context.TODO())
	log.Println("start to watch for ID:", *nodeID)
	for i := 0; i < 10; i++ {
		resp := <-ch
		log.Println("leader changed to", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
	}
}

// ==========Etcd 测试查询命令（query、rev）==========
// 查询主的信息
func query(e1 *concurrency.Election, electName string) {
	// 调用Leader返回主的信息，包括key和value等信息
	resp, err := e1.Leader(context.Background())
	if err != nil {
		log.Printf("failed to get the current leader: %v", err)
	}
	log.Println("current leader:", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
}

// 可以直接查询主的rev信息
func rev(e1 *concurrency.Election, electName string) {
	rev := e1.Rev()
	log.Println("current rev:", rev)
}

// ==========Etcd 测试 Campaign、Proclaim、Resign==========
var count int

// 选主
func elect(e1 *concurrency.Election, electName string) {
	log.Println("acampaigning for ID:", *nodeID)
	// 调用Campaign方法选主,主的值为value-<主节点ID>-<count>
	if err := e1.Campaign(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("campaigned for ID:", *nodeID)
	count++
}

// 为主设置新值
func proclaim(e1 *concurrency.Election, electName string) {
	log.Println("proclaiming for ID:", *nodeID)
	// 调用Proclaim方法设置新值,新值为value-<主节点ID>-<count>
	if err := e1.Proclaim(context.Background(), fmt.Sprintf("value-%d-%d", *nodeID, count)); err != nil {
		log.Println(err)
	}
	log.Println("proclaimed for ID:", *nodeID)
	count++
}

// 重新选主，有可能另外一个节点被选为了主
func resign(e1 *concurrency.Election, electName string) {
	log.Println("resigning for ID:", *nodeID)
	// 调用Resign重新选主
	if err := e1.Resign(context.TODO()); err != nil {
		log.Println(err)
	}
	log.Println("resigned for ID:", *nodeID)
}

// ==========Etcd 测试==========
var ( // 可以设置一些参数，比如节点ID
	nodeID    = flag.Int("id", 0, "node ID")
	addr      = flag.String("addr", "http://127.0.0.1:2379", "etcd addresses")
	electName = flag.String("name", "my-test-elect", "election name")
)

func EtcdTestDemo() {
	flag.Parse()
	// 将etcd的地址解析成slice of string
	endpoints := strings.Split(*addr, ",")
	// 生成一个etcd的clien
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// 创建session,如果程序宕机导致session断掉，etcd能检测到
	session, err := concurrency.NewSession(cli)
	defer session.Close()
	// 生成一个选举对象。下面主要使用它进行选举和查询等操作
	// 另一个方法ResumeElection可以使用既有的leader初始化Election
	//e1 := concurrency.NewElection(session, *electName)
	e1 := concurrency.NewElection(cli, *electName)
	// 从命令行读取命令
	consolescanner := bufio.NewScanner(os.Stdin)
	for consolescanner.Scan() {
		action := consolescanner.Text()
		switch action {
		case "elect": // 选举命令
			go elect(e1, *electName)
		case "proclaim": // 只更新leader的value
			proclaim(e1, *electName)
		case "resign": // 辞去leader,重新选举
			resign(e1, *electName)
		case "watch": // 监控leader的变动
			go watch(e1, *electName)
		case "query": // 查询当前的leader
			query(e1, *electName)
		case "rev":
			rev(e1, *electName)
		default:
			fmt.Println("unknown action")
		}
	}
}
