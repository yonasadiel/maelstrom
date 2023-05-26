package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type EchoReq struct {
	Type  string `json:"type"`
	MsgId int64  `json:"msg_id"`
	Echo  string `json:"echo"`
}

type EchoResp struct {
	Type  string `json:"type"`
	MsgId int64  `json:"msg_id"`
	Echo  string `json:"echo"`
}

func Echo(msg maelstrom.Message) (any, error) {
	reqBody := EchoReq{}
	if err := json.Unmarshal(msg.Body, &reqBody); err != nil {
		return nil, err
	}

	respBody := EchoResp{
		Type:  "echo_ok",
		MsgId: reqBody.MsgId,
		Echo:  reqBody.Echo,
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
	n := maelstrom.NewNode()
	n.Handle("echo", wrapHandler(n, Echo))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
