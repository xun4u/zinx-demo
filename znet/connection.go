package znet

import (
	"errors"
	"fmt"
	"github.com/xun4u/zinx-demo/zinface"
	"io"
	"net"
)

//当前链接模块

type Connection struct {
	//当前链接的socket
	Conn *net.TCPConn
	//链接的ID
	ConnID uint32
	//当前链接状态
	isClosed bool
	//当前链接锁绑定的业务处理方法API
	//handleApi zinface.HandleFunc
	//告知当前链接已经退出/停止的 channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router zinface.IRouter
}

//初始化
func NewConnection(conn *net.TCPConn, connID uint32, router zinface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
	return c
}

//链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("链接读的协程开启")
	defer fmt.Println("ConnID=", c.ConnID, "读取结束，远程地址是", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中，最大512
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("buf读取错误", err)
		//	continue
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端msg head二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err:", err)
		}
		//拆包，得到msgId和len 放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack err", err)
			break
		}
		//根据len再次读取data，放在msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err", err)
				break
			}
		}
		msg.SetData(data)

		////调用当前链接绑定的API处理业务
		//if err := c.handleApi(c.Conn, buf, n); err != nil {
		//	fmt.Println("ConnID=", c.ConnID, "handle错误", err)
		//	break
		//}

		//得到当前conn数据的request
		req := Request{
			conn: c,
			msg:  msg,
		}
		//执行注册的路由方法
		go func(request zinface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		//从路由中找到注册绑定的conn对应的router调用

	}
}

//启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("链接开启，ConnID", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	//todo 启动写数据的业务
}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("链接关闭，ConnID", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket链接
	c.Conn.Close()

	close(c.ExitChan)
}

//获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端tcp状态 ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//提供一个SendMsg方法，将我们要发送给客户端的数据，先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("链接已经关闭，不可发送")
	}

	//封包
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessagePackage(msgId, data))
	if err != nil {
		fmt.Println("pack err msg id=", msgId)
		return errors.New("pack err msg")
	}
	//发送
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("write msg id", msgId, "err:", err)
		return errors.New("conn write err")
	}
	return nil
}
