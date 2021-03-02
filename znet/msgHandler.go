package znet

import (
	"fmt"
	"github.com/xun4u/zinx-demo/utils"
	"github.com/xun4u/zinx-demo/zinface"
	"strconv"
)

//消息处理模块

type MsgHandle struct {
	//存放每个msgid 对应的处理方法
	Apis map[uint32]zinface.IRouter
	//负责worker取任务的消息队列
	TaskQueue []chan zinface.IRequest
	//业务工作worker池的worker数量  这个数量和上面的消息队列是一一对应的
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]zinface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan zinface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

//调度，执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request zinface.IRequest) {
	//从request中找到msgid
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgid=", request.GetMsgID(), "is not found,need register")
		return
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

//启动一个worker工作池(只能发生一次)
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize分别开启Worker，每个worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1当前的worker对应的channel
		mh.TaskQueue[i] = make(chan zinface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//2启动当前worker 阻塞等待消息从channel传进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

//启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan zinface.IRequest) {
	fmt.Println("workerID=", workerID, "is started")
	//不断阻塞等待对应消息队列消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的request，执行当前request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给taskQueue 由worker处理
func (mh MsgHandle) SendMsgToTaskQueue(request zinface.IRequest) {
	//1 将消息平均分配给不同的worker
	//根据客户端建立的connID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("add ConnID=", request.GetConnection().GetConnID(), "request msgid=", request.GetMsgID(), "to workerid=", workerID)

	//2 将消息发送给对应的worker的TaskQueue
	mh.TaskQueue[workerID] <- request
}
