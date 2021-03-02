package zinface

//把客户端请求的链接信息和请求信息包装到request中

type IRequest interface {
	//得到当前链接
	GetConnection() IConnection
	//得到请求的消息数据
	GetData() []byte
}
