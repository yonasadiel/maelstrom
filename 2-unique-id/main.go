package main

import (
	"log"
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

type Server struct {
	id int64
}

func (s *Server) UniqueIds(msg maelstrom.Message) (any, error) {
	newId := atomic.AddInt64(&s.id, 1)
	respBody := UniqueIdsResp{
		Type: "generate_ok",
		Id:   msg.Dest + "-" + strconv.FormatInt(newId, 10),
	}
	return respBody, nil
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
	n.Handle("generate", wrapHandler(n, s.UniqueIds))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
