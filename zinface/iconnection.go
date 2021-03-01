package zinface

import "net"

//定义链接的抽象层
type IConnection interface {
	//启动链接 让当前的链接准备开始工作
	Start()
	//停止链接 结束当前链接的工作
	Stop()
	//获取当前链接的绑定socket conn
	GetTCPConnection() *net.TCPConn
	//获取当前链接模块的链接id
	GetConnID() uint32
	//获取远程客户端tcp状态 ip port
	RemoteAddr() net.Addr
	//发送数据 将数据发送给远程客户端,即往socket中write字节
	Send(data []byte) error
}

//定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error