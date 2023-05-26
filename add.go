package main

import (
	"context"
	"encoding/json"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	KVKeyLock    = "lock"
	KVKeyCounter = "counter"
)

type AddReq struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}

type AddResp struct {
	Type string `json:"type"`
}

func (s *Server) Add(msg maelstrom.Message) (any, error) {
	req := AddReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}

	s.atomic(KVKeyCounter, func(val any) any {
		return val.(int) + req.Delta
	})

	return AddResp{Type: MsgTypeAddOk}, nil
}

func (s *Server) atomic(key string, f func(val any) any) error {
	ctx := context.Background()
	for {
		val, err := s.seqKV.Read(ctx, key)
		if rpcErr, ok := err.(*maelstrom.RPCError); ok && rpcErr.Code == maelstrom.KeyDoesNotExist {
			val = nil
		} else if err != nil {
			return err
		}

		newVal := f(val)
		if newVal == val {
			return nil
		}

		err = s.seqKV.CompareAndSwap(ctx, key, val, newVal, true)
		if rpcErr, ok := err.(*maelstrom.RPCError); ok && rpcErr.Code == maelstrom.PreconditionFailed {
			//
		} else if err != nil {
			return err
		} else {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
}
