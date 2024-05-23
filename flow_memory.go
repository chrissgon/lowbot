package lowbot

import (
	"fmt"
)

type MemoryFlowPersist struct {
	Sessions map[string]*Flow
}

func NewMemoryFlowPersist() (FlowPersist, error) {
	loc := &MemoryFlowPersist{Sessions: map[string]*Flow{}}

	return loc, nil
}

func (memory *MemoryFlowPersist) Get(sessionID string) (*Flow, error) {
	flow := memory.Sessions[sessionID]

	if flow == nil {
		return nil, fmt.Errorf("not found flow")
	}

	return flow, nil
}

func (memory *MemoryFlowPersist) Set(flow *Flow) error {
	memory.Sessions[flow.SessionID] = flow

	return nil
}
