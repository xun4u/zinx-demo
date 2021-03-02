package zinface

//将请求的消息封装到message中
type IMessage interface {
	//获取消息
	GetMsgId() uint32
	GetMsgLen() uint32
	GetMsgData() []byte

	//设置消息
	SetMsgId(uint32)
	SetDataLen(uint32)
	SetData([]byte)
}
