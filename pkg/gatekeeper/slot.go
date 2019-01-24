package gatekeeper

type AgentSlot struct {
	Host    string
	Port    uint16
	AgentID string
}

type GateKeeper struct {
}

func (g *GateKeeper) AllocateAgentSlot() {

}
