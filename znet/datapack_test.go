package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只负责测试datapack拆包封包的单元测试
func TestDataPack(t *testing.T) {

	//模拟服务器
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	//创建一个go承载负责客户端业务处理
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}

			go func(conn net.Conn) {
				//处理客户端的请求
				//拆包过程
				dp := NewDataPack()
				for {
					//第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err")
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//第二次从conn读，根据head中中的len来读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据len长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						//完整的一个消息已经读取完毕
						fmt.Println("msgId=", msg.Id, "msgLen=", msg.DataLen, "msgData=", string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
	}

	//创建一个封包对象
	dp := NewDataPack()

	//模拟粘包过程，封装2个msg一同发送
	//第一个msg
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
	}
	//第二个msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
	}
	//2个包黏在一起发送
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}
}
