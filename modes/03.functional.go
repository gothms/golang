package modes

import (
	"crypto/tls"
	"time"
)

/*
Go 编程模式：Functional Options

函数式编程
	目前 Go 语言中最流行的一种编程模式

配置选项问题
	对一个对象（或是业务实体）进行相关的配置
	要有侦听的 IP 地址 Addr 和端口号 Port ，这两个配置选项是必填的（当然，IP 地址和端口号都可以有默认值，不过这里我们用于举例，所以是没有默认值，而且不能为空，需要是必填的）
	然后，还有协议 Protocol 、 Timeout 和MaxConns 字段，这几个字段是不能为空的，但是有默认值的，比如，协议是 TCP，超时30秒 和 最大链接数1024个
	还有一个 TLS ，这个是安全链接，需要配置相关的证书和私钥。这个是可以为空的
配置对象方案
	最常见的方式是使用一个配置对象
		type Config struct {
			Protocol string
			Timeout  time.Duration
			Maxconns int
			TLS      *tls.Config
		}
	把那些非必输的选项都移到一个结构体里
		type Server struct {
			Addr string
			Port int
			Conf *Config
		}
	于是，我们就只需要一个 NewServer() 的函数了，在使用前需要构造 Config 对象
		func NewServer(addr string, port int, conf *Config) (*Server, error) {
			//...
		}

		//Using the default configuratrion
		srv1, _ := NewServer("localhost", 9000, nil)

		conf := ServerConfig{Protocol:"tcp", Timeout: 60*time.Duration}
		srv2, _ := NewServer("locahost", 9000, &conf)

Builder 模式
	Java Builder 模式
		User user = new User.Builder()
		  .name("Hao Chen")
		  .email("haoel@hotmail.com")
		  .nickname("左耳朵")
		  .build();
	仿照这个模式
		//使用一个builder类来做包装
		type ServerBuilder struct {
		  Server
		}

		func (sb *ServerBuilder) Create(addr string, port int) *ServerBuilder {
		  sb.Server.Addr = addr
		  sb.Server.Port = port
		  //其它代码设置其它成员的默认值
		  return sb
		}

		func (sb *ServerBuilder) WithProtocol(protocol string) *ServerBuilder {
		  sb.Server.Protocol = protocol
		  return sb
		}

		func (sb *ServerBuilder) WithMaxConn( maxconn int) *ServerBuilder {
		  sb.Server.MaxConns = maxconn
		  return sb
		}

		func (sb *ServerBuilder) WithTimeOut( timeout time.Duration) *ServerBuilder {
		  sb.Server.Timeout = timeout
		  return sb
		}

		func (sb *ServerBuilder) WithTLS( tls *tls.Config) *ServerBuilder {
		  sb.Server.TLS = tls
		  return sb
		}

		func (sb *ServerBuilder) Build() (Server) {
		  return  sb.Server
		}
	使用方式
		sb := ServerBuilder{}
		server, err := sb.Create("127.0.0.1", 8080).
		  WithProtocol("udp").
		  WithMaxConn(1024).
		  WithTimeOut(30*time.Second).
		  Build()
	想省掉这个包装的结构体，就要请出 Functional Options 上场了：函数式编程

Functional Options
	示例代码如下
		定义一个函数类型
		使用函数式的方式定义一组的函数
		再定一个 NewServer()的函数
	高度整洁和优雅
		不但解决了“使用 Config 对象方式的需要有一个 config 参数，但在不需要的时候，是放 nil 还是放 Config{}”的选择困难问题
		也不需要引用一个 Builder 的控制对象，直接使用函数式编程，在代码阅读上也很优雅
	Functional Options 至少带来 6 个好处
		直觉式的编程
		高度的可配置化
		很容易维护和扩展
		自文档
		新来的人很容易上手
		没有什么令人困惑的事（是 nil 还是空）
	参考
		http://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html
*/

type Server struct {
	Addr     string
	Port     int
	Protocol string
	Timeout  time.Duration
	MaxConns int
	TLS      *tls.Config
}
type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}
func Timeout(t time.Duration) Option {
	return func(s *Server) {
		s.Timeout = t
	}
}
func MaxConns(mc int) Option {
	return func(s *Server) {
		s.MaxConns = mc
	}
}
func TLS(tls *tls.Config) Option {
	return func(s *Server) {
		s.TLS = tls
	}
}
func NewServer(addr string, port int, options ...func(*Server)) (*Server, error) {
	srv := &Server{
		Addr:     addr,
		Port:     port,
		Protocol: "tcp",
		Timeout:  30 * time.Second,
		MaxConns: 1000,
		TLS:      nil,
	}
	for _, option := range options {
		option(srv)
	}
	// ...
	return srv, nil
}
func test() {
	//s1, _ := NewServer("localhost", 1024)
	//s2, _ := NewServer("localhost", 2048, Protocol("udp"))
	//s3, _ := NewServer("0.0.0.0", 8080, Timeout(300*time.Second), MaxConns(1000))
}
