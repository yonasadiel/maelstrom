package main

import (
	"log"
	"os"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	MsgTypeInit        = "init"
	MsgTypeInitOk      = "init_ok"
	MsgTypeEcho        = "echo"
	MsgTypeEchoOk      = "echo_ok"
	MsgTypeGenerate    = "generate"
	MsgTypeGenerateOk  = "generate_ok"
	MsgTypeBroadcast   = "broadcast"
	MsgTypeBroadcastOk = "broadcast_ok"
	MsgTypeRead        = "read"
	MsgTypeReadOk      = "read_ok"
	MsgTypeTopology    = "topology"
	MsgTypeTopologyOk  = "topology_ok"
	MsgTypeAdd         = "add"
	MsgTypeAddOk       = "add_ok"

	WorkloadBroadcast       = "broadcast"
	WorkloadGrowOnlyCounter = "grow-only-counter"
)

type Server struct {
	n           *maelstrom.Node
	seqKV       *maelstrom.KV
	initialized chan struct{}
	workload    string

	// for Unique ID module
	id int64

	// for Broadcast module
	broadcastedLock sync.RWMutex
	broadcasted     []int64
	broadcastedSet  map[int64]struct{}
}

func (s *Server) wrapHandler(f func(msg maelstrom.Message) (any, error), waitForInit bool) func(msg maelstrom.Message) error {
	return func(msg maelstrom.Message) error {
		if waitForInit {
			<-s.initialized
		}
		resp, err := f(msg)
		if err != nil {
			return err
		}
		return s.n.Reply(msg, resp)
	}
}

func main() {
	n := maelstrom.NewNode()
	seqKV := maelstrom.NewSeqKV(n)
	s := Server{
		n:               n,
		seqKV:           seqKV,
		initialized:     make(chan struct{}),
		workload:        os.Getenv("MWORKLOAD"),
		id:              0,
		broadcastedLock: sync.RWMutex{},
		broadcasted:     make([]int64, 0),
		broadcastedSet:  make(map[int64]struct{}, 0),
	}
	n.Handle("init", s.wrapHandler(s.Init, false))
	n.Handle("echo", s.wrapHandler(s.Echo, false))
	n.Handle("generate", s.wrapHandler(s.UniqueIds, false))
	n.Handle("broadcast", s.wrapHandler(s.Broadcast, false))
	n.Handle("read", s.wrapHandler(s.Read, true))
	n.Handle("topology", s.wrapHandler(s.Topology, false))
	n.Handle("add", s.wrapHandler(s.Add, true))

	if s.workload == WorkloadBroadcast {
		go s.sendPendingBroadcast()
	}

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
