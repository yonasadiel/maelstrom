package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ReadReq struct {
	Type string `json:"type"`
}

type ReadResp struct {
	Type     string  `json:"type"`
	Messages []int64 `json:"messages"`
}

func (s *Server) Read(msg maelstrom.Message) (any, error) {
	req := ReadReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}
	resp := ReadResp{Type: MsgTypeReadOk}
	s.broadcastedLock.RLock()
	resp.Messages = make([]int64, len(s.broadcasted))
	for i := range s.broadcasted {
		resp.Messages[i] = s.broadcasted[i]
	}
	s.broadcastedLock.RUnlock()
	return resp, nil
}
