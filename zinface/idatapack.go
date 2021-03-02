package zinface

//封包 拆包模块，直接面向TCP链接中的数据流，处理粘包问题

type IDataPack interface {
	//获取包头长度方法
	GetHeadLen() uint32
	//封包方法
	Pack(IMessage) ([]byte, error)
	//拆包方法
	UnPack([]byte) (IMessage, error)
}
