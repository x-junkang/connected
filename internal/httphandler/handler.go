package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/x-junkang/connected/pkg/ciface"
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
			log.Err(err).Msg("json marshal fail")
			return
		}
		rw.Write(respData)
	})
}

type AllConnResp struct {
	Count int `json:"len"`
}
