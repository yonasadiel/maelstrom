package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type SendReq struct {
	Type string `json:"type"`
	Key  string `json:"key"`
	Msg  int64  `json:"msg"`
}

type SendResp struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
}

type KafkaMessage struct {
	Message int64
	Offset  int64
}

func (s *Server) Send(msg maelstrom.Message) (any, error) {
	req := SendReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}

	s.msgsLock.Lock()
	defer s.msgsLock.Unlock()

	if _, ok := s.msgs[req.Key]; !ok {
		s.msgs[req.Key] = make([]KafkaMessage, 0)
	}
	lastOffset := int64(-1)
	existingMessagesNum := len(s.msgs[req.Key])
	if existingMessagesNum > 0 {
		lastOffset = s.msgs[req.Key][existingMessagesNum-1].Offset
	}
	newMsg := KafkaMessage{
		Message: req.Msg,
		Offset:  lastOffset + 1,
	}
	s.msgs[req.Key] = append(s.msgs[req.Key], newMsg)

	resp := SendResp{
		Type:   MsgTypeSendOk,
		Offset: newMsg.Offset,
	}
	return resp, nil
}
