package main

import (
	"encoding/json"
	"math/rand"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/net/context"
)

const (
	GossipPeriod  = 30 * time.Millisecond
	GossipTimeout = 300 * time.Millisecond
)

type BroadcastReq struct {
	Type    string `json:"type"`
	Message *int64 `json:"message,omitempty"`
	// Added by myself, for gossip protocol
	Messages []int64 `json:"messages,omitempty"`
}

type BroadcastResp struct {
	Type string `json:"type"`
}

func (s *Server) Broadcast(msg maelstrom.Message) (any, error) {
	req := BroadcastReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}

	missing := make([]int64, 0)
	incoming := req.Messages[:]
	if req.Message != nil {
		incoming = append(incoming, *req.Message)
	}
	s.broadcastedLock.RLock()
	for _, m := range incoming {
		if _, ok := s.broadcastedSet[m]; !ok {
			missing = append(missing, m)
		}
	}
	s.broadcastedLock.RUnlock()

	if len(missing) > 0 {
		s.broadcastedLock.Lock()
		for _, m := range missing {
			s.broadcasted = append(s.broadcasted, m)
			s.broadcastedSet[m] = struct{}{}
		}
		s.broadcastedLock.Unlock()
	}

	resp := BroadcastResp{Type: MsgTypeBroadcastOk}
	return resp, nil
}

func (s *Server) sendPendingBroadcast() {
	ticker := time.NewTicker(GossipPeriod)
	for range ticker.C {
		nodeIDs := s.n.NodeIDs()
		if len(nodeIDs) == 0 {
			continue
		}

		s.broadcastedLock.RLock()
		messages := make([]int64, len(s.broadcasted))
		for i := range s.broadcasted {
			messages[i] = s.broadcasted[i]
		}
		s.broadcastedLock.RUnlock()

		req := BroadcastReq{
			Type:     MsgTypeBroadcast,
			Messages: messages,
		}
		dest := nodeIDs[rand.Intn(len(nodeIDs))]
		go func() {
			ctx, cancelFn := context.WithTimeout(context.Background(), GossipTimeout)
			defer cancelFn()
			if _, err := s.n.SyncRPC(ctx, dest, req); err != nil {
				// TODO handle error
			}
		}()
	}
}
