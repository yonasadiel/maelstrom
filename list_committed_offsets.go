package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ListCommittedOffsetsReq struct {
	Type string   `json:"type"`
	Keys []string `json:"keys"`
}

type ListCommittedOffsetsResp struct {
	Type    string           `json:"type"`
	Offsets map[string]int64 `json:"offsets"`
}

func (s *Server) ListCommittedOffsets(msg maelstrom.Message) (any, error) {
	req := ListCommittedOffsetsReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}

	s.offsetsLock.RLock()
	defer s.offsetsLock.RUnlock()

	resp := ListCommittedOffsetsResp{
		Type:    MsgTypeListCommittedOffsetsOk,
		Offsets: make(map[string]int64, len(req.Keys)),
	}
	for _, key := range req.Keys {
		resp.Offsets[key] = s.offsets[key]
	}
	return resp, nil
}
