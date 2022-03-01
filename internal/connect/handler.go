package connect

import "github.com/x-junkang/connected/pkg/ciface"

type TcpHandler struct {
	Apis map[uint32]ciface.IRouter
}

func NewTcpHandler() *TcpHandler {
	return &TcpHandler{
		Apis: make(map[uint32]ciface.IRouter, 10),
	}
}

func (handler *TcpHandler) DoMsgHandler(req ciface.IRequest) {
	id := req.GetMsgID()
	if router, ok := handler.Apis[id]; ok {
		router.PreHandle(req)
		router.Handle(req)
		router.PreHandle(req)
	} else {
		// default handle
	}
}

func (handler *TcpHandler) AddRouter(msgID uint32, router ciface.IRouter) {
	handler.Apis[msgID] = router
}
