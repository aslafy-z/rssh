package gatekeeper

import "errors"

type AgentSlot struct {
	Host    string `json:"host"`
	Port    uint16 `json:"port"`
	AgentID string `json:"agentID"`
}

type Gate struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

type GateKeeper struct {
	backends []Gate
}

func NewGateKeeper(addr string, port uint16) (*GateKeeper, error) {
	return nil, errors.New("NotImplemented")
}

func (g *GateKeeper) AllocateAgentSlot() (*AgentSlot, error) {
	return nil, errors.New("NotImplemented")
}
