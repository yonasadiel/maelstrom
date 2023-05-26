package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type EchoReq struct {
	Type  string `json:"type"`
	MsgId int64  `json:"msg_id"`
	Echo  string `json:"echo"`
}

type EchoResp struct {
	Type  string `json:"type"`
	MsgId int64  `json:"msg_id"`
	Echo  string `json:"echo"`
}

func (s *Server) Echo(msg maelstrom.Message) (any, error) {
	reqBody := EchoReq{}
	if err := json.Unmarshal(msg.Body, &reqBody); err != nil {
		return nil, err
	}

	respBody := EchoResp{
		Type:  MsgTypeEchoOk,
		MsgId: reqBody.MsgId,
		Echo:  reqBody.Echo,
	}
	return respBody, nil
}
