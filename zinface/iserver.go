package zinface

//定义一个服务器接口
type IServer interface {

	//启动
	Start()
	//停止
	Stop()
	//运行
	Serve()

	//路由功能，给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)

	GetConnManager() IConnManager

	//注册OnConnStart
	SetOnConnStart(func(connection IConnection))
	//注册OnConnStop
	SetOnConnStop(func(connection IConnection))
	//调用OnConnStart
	CallOnConnStart(connection IConnection)
	//调用OnConnStop
	CallOnConnStop(connection IConnection)
}
