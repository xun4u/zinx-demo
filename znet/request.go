package znet

import "github.com/xun4u/zinx-demo/zinface"

type Request struct {
	//已经和客户端建立好的链接
	conn zinface.IConnection

	//客户端请求的数据
	msg zinface.IMessage
}

//得到当前链接
func (r *Request) GetConnection() zinface.IConnection {
	return r.conn
}

//得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetMsgData()
}

//得到msgid
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
