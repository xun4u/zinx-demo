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
}

//初始化方法
func NewServer(name string) zinface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}

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

		//3 accept 阻塞客户端的链接 处理客户端的链接业务（读写）
		for {
			//如果有客户端连接进来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("客户端连接失败：", err)
				continue
			}

			//已经建立了连接，然后业务操作(先读sock 然后写入sock)
			go func() {
				for {
					buf := make([]byte, 512)
					n, err := conn.Read(buf)
					if err != nil {
						fmt.Println("读取buf错误：", err)
						continue
					}

					//回显
					if _, err := conn.Write(buf[:n]); err != nil {
						fmt.Println("回写buf错误：", err)
						continue
					}
				}
			}()
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
