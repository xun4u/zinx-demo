package utils

import (
	"encoding/json"
	"github.com/xun4u/zinx-demo/zinface"
	"io/ioutil"
)

//存储一切有关zinx框架的全局参数，供其他模块使用
//一些参数是可以通过zinx.json由用户进行配置

type GlobalObj struct {
	//server
	TcpServer zinface.IServer //当前zinx全局server对象
	Host      string          //当前服务器主机监听ip
	TcpPort   int             //当前服务器监听端口
	Name      string          //服务器名称

	//zinx
	Version        string //zinx版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前zinx框架数据包的最大值
}

//定义一个全局对象
var GlobalObject *GlobalObj

//从zinx.json 去加载自定义配置
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	//将json文件数据解析到结构体中
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

//初始化global
func init() {
	//如果配置文件没有加载默认的值
	GlobalObject = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "zinx",
		Version:        "v0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	GlobalObject.Reload()
}
