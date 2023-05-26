package main

import (
	"context"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type InitReq struct {
	Type string `json:"type"`
}

type InitResp struct {
	Type string `json:"type"`
}

func (s *Server) Init(msg maelstrom.Message) (any, error) {
	// req := InitReq{}
	// if err := json.Unmarshal(msg.Body, &req); err != nil {
	// 	return nil, err
	// }

	if s.workload == WorkloadGrowOnlyCounter {
		if err := s.seqKV.Write(context.Background(), KVKeyCounter, 0); err != nil {
			return nil, err
		}
	}

	close(s.initialized)
	resp := InitResp{Type: MsgTypeInitOk}
	return resp, nil
}
