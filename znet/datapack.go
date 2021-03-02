package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/xun4u/zinx-demo/utils"
	"github.com/xun4u/zinx-demo/zinface"
)

type DataPack struct{}

//初始化
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	//datalen uint32（4字节） + ID uint32
	return 8
}

//封包方法
func (dp *DataPack) Pack(msg zinface.IMessage) ([]byte, error) {
	//创建一个存放byte字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将len写入buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将id写入buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data写入buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

//拆包方法 将包的head信息读出来，根据head里的信息data的长度再进行一次读
func (dp *DataPack) UnPack(binaryData []byte) (zinface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}
	//只解压head信息

	//读len
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断datalen是否超出了允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("收到的包超过限制")
	}

	return msg, nil
}
