package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type CommitOffsetsReq struct {
	Type    string           `json:"type"`
	Offsets map[string]int64 `json:"offsets"`
}

type CommitOffsetsResp struct {
	Type string `json:"type"`
}

func (s *Server) CommitOffsets(msg maelstrom.Message) (any, error) {
	req := CommitOffsetsReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}

	s.offsetsLock.Lock()
	defer s.offsetsLock.Unlock()

	for key, offset := range req.Offsets {
		s.offsets[key] = offset
	}

	resp := CommitOffsetsResp{Type: MsgTypeCommitOffsetsOk}
	return resp, nil
}
