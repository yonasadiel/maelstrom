package main

import (
	"context"
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ReadReq struct {
	Type string `json:"type"`
}

type ReadResp struct {
	Type     string   `json:"type"`
	Value    *int     `json:"value,omitempty"`    // for grow-only-counter
	Messages *[]int64 `json:"messages,omitempty"` // for broadcast
}

func (s *Server) Read(msg maelstrom.Message) (any, error) {
	req := ReadReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}
	resp := ReadResp{Type: MsgTypeReadOk}

	var err error
	switch s.workload {
	case WorkloadBroadcast:
		err = s.readBroadcast(req, &resp)
	case WorkloadGrowOnlyCounter:
		err = s.readGrowOnlyCounter(req, &resp)
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) readBroadcast(req ReadReq, resp *ReadResp) error {
	s.broadcastedLock.RLock()
	defer s.broadcastedLock.RUnlock()

	messages := make([]int64, len(s.broadcasted))
	for i := range s.broadcasted {
		messages[i] = s.broadcasted[i]
	}
	resp.Messages = &messages
	return nil
}

func (s *Server) readGrowOnlyCounter(req ReadReq, resp *ReadResp) error {
	ctx := context.Background()
	val, err := s.seqKV.ReadInt(ctx, KVKeyCounter)
	if err != nil {
		return err
	}
	resp.Value = &val
	return nil
}
