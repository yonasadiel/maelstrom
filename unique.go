package main

import (
	"strconv"
	"sync/atomic"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type UniqueIdsReq struct {
	Type string `json:"type"`
}

type UniqueIdsResp struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

func (s *Server) UniqueIds(msg maelstrom.Message) (any, error) {
	newId := atomic.AddInt64(&s.id, 1)
	resp := UniqueIdsResp{
		Type: MsgTypeGenerateOk,
		Id:   msg.Dest + "-" + strconv.FormatInt(newId, 10),
	}
	return resp, nil
}
