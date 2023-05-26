package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastReq struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
}

type BroadcastResp struct {
	Type string `json:"type"`
}

type ReadReq struct {
	Type string `json:"type"`
}

type ReadResp struct {
	Type     string  `json:"type"`
	Messages []int64 `json:"messages"`
}

type TopologyReq struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type TopologyResp struct {
	Type string `json:"type"`
}

type Server struct {
	dataLock sync.RWMutex
	data     []int64
}

func (s *Server) Broadcast(msg maelstrom.Message) (any, error) {
	req := BroadcastReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}
	s.dataLock.Lock()
	s.data = append(s.data, req.Message)
	s.dataLock.Unlock()
	resp := BroadcastResp{Type: "broadcast_ok"}
	return resp, nil
}

func (s *Server) Read(msg maelstrom.Message) (any, error) {
	req := ReadReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}
	resp := ReadResp{Type: "read_ok"}
	s.dataLock.RLock()
	resp.Messages = make([]int64, len(s.data))
	for i := range s.data {
		resp.Messages[i] = s.data[i]
	}
	s.dataLock.RUnlock()
	return resp, nil
}

func (s *Server) Topology(msg maelstrom.Message) (any, error) {
	req := TopologyReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}
	// TODO
	resp := TopologyResp{Type: "topology_ok"}
	return resp, nil
}

func wrapHandler(n *maelstrom.Node, f func(msg maelstrom.Message) (any, error)) func(msg maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		resp, err := f(msg)
		if err != nil {
			return err
		}
		return n.Reply(msg, resp)
	}
}

func main() {
	s := Server{}
	n := maelstrom.NewNode()
	n.Handle("broadcast", wrapHandler(n, s.Broadcast))
	n.Handle("read", wrapHandler(n, s.Read))
	n.Handle("topology", wrapHandler(n, s.Topology))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
