package znet

type Message struct {
	Id      uint32 //消息id
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

func NewMessagePackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

//获取消息
func (m *Message) GetMsgId() uint32 {
	return m.Id
}
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}
func (m *Message) GetMsgData() []byte {
	return m.Data
}

//设置消息
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
