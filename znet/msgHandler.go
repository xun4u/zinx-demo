package znet

import (
	"fmt"
	"github.com/xun4u/zinx-demo/zinface"
	"strconv"
)

//消息处理模块

type MsgHandle struct {
	//存放每个msgid 对应的处理方法
	Apis map[uint32]zinface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{Apis: make(map[uint32]zinface.IRouter)}
}

//调度，执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request zinface.IRequest) {
	//从request中找到msgid
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgid=", request.GetMsgID(), "is not found,need register")
	}
	//根据msgid调度对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router zinface.IRouter) {
	//1判断当前的msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api,msgid=" + strconv.Itoa(int(msgID)))
	}
	//2添加msg和API绑定关系
	mh.Apis[msgID] = router
	fmt.Println("add api msgid-router成功", msgID)
}
