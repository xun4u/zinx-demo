package znet

import "github.com/xun4u/zinx-demo/zinface"

//实现router时先嵌入这个base基类，然后根据需要对方法进行重写
type BaseRouter struct{}

//之所以下面不实现，是因为有的router不一定全部要实现以下3个业务，所以router全部继承baserouter的好处就可以按需重写

//在处理conn之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request zinface.IRequest) {}

//在处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request zinface.IRequest) {}

//在处理conn之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request zinface.IRequest) {}
