package znet

import (
	"fmt"
	"github.com/xun4u/zinx-demo/zinface"
	"net"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//名称
	Name string
	//绑定的ip版本
	IPVersion string
	//监听的ip
	IP string
	//监听端口
	Port int
	//当前的server添加一个router，server注册的链接对应的处理业务
	Router zinface.IRouter
}

//初始化方法
func NewServer(name string) zinface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}

//定义当前客户连接锁绑定的handle 目前是写死 以后根据用户自定义
//func CallBackToClient(conn *net.TCPConn, data []byte, n int) error {
//	//回显的业务
//	fmt.Println("链接handle调用CallBackToClient")
//	if _, err := conn.Write(data[:n]); err != nil {
//		fmt.Println("回写buf错误：", err)
//		return errors.New("CallBackToClient err")
//	}
//	return nil
//}

func (s *Server) Start() {

	go func() {
		//1 tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("获取tcp地址失败：", err)
			return
		}
		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("监听失败：", err)
			return
		}
		fmt.Println("服务器开启成功")
		var cid uint32
		cid = 0

		//3 accept 阻塞客户端的链接 处理客户端的链接业务（读写）
		for {
			//如果有客户端连接进来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("客户端连接失败：", err)
				continue
			}

			//将处理新链接业务的方法 和conn绑定 得到链接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			//启动当前链接业务处理
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	//todo 将一些服务器资源，状态或者一些已经开辟的链接信息 进行停止或回收
}

func (s *Server) Serve() {
	s.Start()
	//todo 以后做一些服务器启动后额外的业务

	select {}
}

func (s *Server) AddRouter(router zinface.IRouter) {
	s.Router = router
	fmt.Println("添加路由成功")
}
