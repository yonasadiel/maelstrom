package main

import maelstrom "github.com/jepsen-io/maelstrom/demo/go"

type TopologyReq struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type TopologyResp struct {
	Type string `json:"type"`
}

func (s *Server) Topology(msg maelstrom.Message) (any, error) {
	// req := TopologyReq{}
	// if err := json.Unmarshal(msg.Body, &req); err != nil {
	// 	return nil, err
	// }
	resp := TopologyResp{Type: "topology_ok"}
	return resp, nil
}
