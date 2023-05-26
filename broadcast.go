package main

import (
	"encoding/json"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/net/context"
)

type BroadcastReq struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
}

type BroadcastResp struct {
	Type string `json:"type"`
}

type pendingBroadcast struct {
	dest string
	req  BroadcastReq
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
			if nodeID == s.n.ID() || nodeID == msg.Src {
				continue
			}
			broadcastReq := BroadcastReq{
				Type:    MsgTypeBroadcast,
				Message: req.Message,
			}
			s.broadcastQueue <- pendingBroadcast{nodeID, broadcastReq}
		}
	}

	resp := BroadcastResp{Type: MsgTypeBroadcastOk}
	return resp, nil
}

func (s *Server) sendPendingBroadcast() {
	for b := range s.broadcastQueue {
		ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
		if _, err := s.n.SyncRPC(ctx, b.dest, b.req); err != nil {
			s.broadcastQueue <- b
		}
		cancelFn()
	}
}
