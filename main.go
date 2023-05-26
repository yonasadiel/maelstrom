package main

import (
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	MsgTypeEcho        = "echo"
	MsgTypeEchoOk      = "echo_ok"
	MsgTypeGenerate    = "generate"
	MsgTypeGenerateOk  = "generate_ok"
	MsgTypeBroadcast   = "broadcast"
	MsgTypeBroadcastOk = "broadcast_ok"
	MsgTypeRead        = "read"
	MsgTypeReadOk      = "read_ok"
	MsgTypeTopology    = "topology"
	MsgTypeTopologyOk  = "topology"
)

type Server struct {
	n *maelstrom.Node

	// for Unique ID module
	id int64

	// for Broadcast module
	broadcastedLock sync.RWMutex
	broadcasted     []int64
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
	n := maelstrom.NewNode()
	s := Server{n: n}
	n.Handle("echo", wrapHandler(n, s.Echo))
	n.Handle("generate", wrapHandler(n, s.UniqueIds))
	n.Handle("broadcast", wrapHandler(n, s.Broadcast))
	n.Handle("read", wrapHandler(n, s.Read))
	n.Handle("topology", wrapHandler(n, s.Topology))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
