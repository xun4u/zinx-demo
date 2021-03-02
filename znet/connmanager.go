package znet

import (
	"errors"
	"fmt"
	"github.com/xun4u/zinx-demo/zinface"
	"sync"
)

//链接管理模块
type ConnManager struct {
	connections map[uint32]zinface.IConnection //管理的连接集合
	connLock    sync.RWMutex                   //保护链接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]zinface.IConnection),
	}
}

//添加链接
func (cm *ConnManager) Add(conn zinface.IConnection) {
	//保护共享资源map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to cm succ conn num=", cm.Len())
}

//删除链接
func (cm *ConnManager) Remove(conn zinface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connection remove succ conn num=", cm.Len())
}

//根据connID获取链接
func (cm *ConnManager) Get(connID uint32) (zinface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//得到当前链接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

//清楚并终止所有链接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
	fmt.Println("clear all connections succ", cm.Len())
}
