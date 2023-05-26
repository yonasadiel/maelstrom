package main

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type PollReq struct {
	Type    string           `json:"type"`
	Offsets map[string]int64 `json:"offsets"`
}

type PollResp struct {
	Type string               `json:"type"`
	Msgs map[string][][]int64 `json:"msgs"`
}

func (s *Server) Poll(msg maelstrom.Message) (any, error) {
	req := PollReq{}
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return nil, err
	}

	s.msgsLock.RLock()
	defer s.msgsLock.RUnlock()

	resp := PollResp{
		Type: MsgTypePollOk,
		Msgs: make(map[string][][]int64),
	}

	for key, offset := range req.Offsets {
		i := 0
		for i < len(s.msgs[key]) && s.msgs[key][i].Offset < offset {
			i++
		}
		keyMessages := make([][]int64, 0, 5)
		for i < len(s.msgs[key]) && len(keyMessages) < 5 {
			keyMessages = append(keyMessages, []int64{s.msgs[key][i].Offset, s.msgs[key][i].Message})
			i++
		}
		resp.Msgs[key] = keyMessages
	}

	return resp, nil
}
