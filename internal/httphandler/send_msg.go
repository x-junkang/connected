package httphandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

type SendMsgReq struct {
	MsgID   int64  `json:"msg_id"`
	MsgFrom uint64 `json:"msg_from"`
	MsgTo   uint64 `json:"msg_to"`
	Content string `json:"content"`
}

func (s *HttpServer) SendMsg() func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {

		req := &SendMsgReq{}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Str("handler", "SendMsg").Msg("read data from req fail")
			WriteErr(rw, 404, err.Error())
			return
		}
		err = json.Unmarshal(data, req)
		if err != nil {
			log.Err(err).Msg("unmarshal req fail")
			WriteErr(rw, 404, err.Error())
			return
		}
		toConn, err := s.tcpserver.GetConnMgr().Get(req.MsgTo)
		if err != nil {
			log.Err(err).Msg("get conn fail")
			WriteErr(rw, 404, err.Error())
			return
		}

		toConn.SendMsg(1, []byte(req.Content))

		resp := ConmmonResp{
			Status: 200,
			Msg:    "success",
		}
		respData, err := json.Marshal(resp)
		if err != nil {
			log.Err(err).Msg("json marshal fail")
			return
		}
		rw.Write(respData)
	}
}

func WriteErr(rw http.ResponseWriter, status int, ErrMsg string) {
	resp := ConmmonResp{
		Status: status,
		Msg:    ErrMsg,
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		log.Err(err).Msg("json marshal fail")
		return
	}
	rw.Write(respData)
}
