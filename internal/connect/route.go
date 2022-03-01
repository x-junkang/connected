package connect

import (
	"fmt"

	"github.com/x-junkang/connected/pkg/ciface"
)

type EmptyRouter struct {
}

func (router *EmptyRouter) PreHandle(request ciface.IRequest) {

} //在处理conn业务之前的钩子方法
func (router *EmptyRouter) Handle(request ciface.IRequest) {

} //处理conn业务的方法
func (router *EmptyRouter) PostHandle(request ciface.IRequest) {

}

type HelloRouter struct {
	EmptyRouter
}

func (router *HelloRouter) Handle(req ciface.IRequest) {
	msgID := req.GetMsgID()
	data := req.GetData()
	fmt.Println("handler msg 1")
	req.GetConnection().SendMsg(msgID, data)
}

type SendMsgRouter struct {
	EmptyRouter
}

func (router *SendMsgRouter) Handle(req ciface.IRequest) {
	msgID := req.GetMsgID()
	data := req.GetData()
	fmt.Println("handler msg 2")
	req.GetConnection().SendMsg(msgID, data)
}
