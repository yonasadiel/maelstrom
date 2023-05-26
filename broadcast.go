package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastReq struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
}

type BroadcastResp struct {
	Type string `json:"type"`
}

func (s *Server) Broadcast(msg maelstrom.Message) (any, error) {
	req := BroadcastReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}

	exist := false
	s.broadcastedLock.RLock()
	for i := 0; !exist && i < len(s.broadcasted); i++ {
		if s.broadcasted[i] == req.Message {
			exist = true
		}
	}
	s.broadcastedLock.RUnlock()

	if !exist {
		s.broadcastedLock.Lock()
		s.broadcasted = append(s.broadcasted, req.Message)
		s.broadcastedLock.Unlock()

		for _, nodeID := range s.n.NodeIDs() {
			if nodeID == s.n.ID() {
				continue
			}
			broadcastReq := BroadcastReq{
				Type:    MsgTypeBroadcast,
				Message: req.Message,
			}
			if err := s.n.RPC(nodeID, broadcastReq, nil); err != nil {
				return nil, err
			}
		}
	}

	resp := BroadcastResp{Type: MsgTypeBroadcastOk}
	return resp, nil
}
