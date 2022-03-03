package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/x-junkang/connected/internal/clog"
	"github.com/x-junkang/connected/pkg/ciface"
	"go.uber.org/zap"
)

type HttpServer struct {
	http.ServeMux
	tcpserver ciface.IServer
}

func NewHttpServer(tcpserver ciface.IServer) *HttpServer {
	ht := &HttpServer{
		http.ServeMux{},
		tcpserver,
	}
	ht.router()
	return ht
}

func (s *HttpServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	s.ServeMux.ServeHTTP(resp, req)
}

func (s *HttpServer) router() {
	s.HandleFunc("/all", func(rw http.ResponseWriter, r *http.Request) {
		data := AllConnResp{
			Count: s.tcpserver.GetConnMgr().Len(),
		}
		respData, err := json.Marshal(data)
		if err != nil {
			clog.Logger.Error("json marshal fail", zap.String("err", err.Error()))
			return
		}
		rw.Write(respData)
	})
}

type AllConnResp struct {
	Count int `json:"len"`
}
