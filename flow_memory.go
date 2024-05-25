package lowbot

import (
	"fmt"
)

type MemoryFlowPersist struct {
	Sessions map[any]*Flow
}

func NewMemoryFlowPersist() (FlowPersist, error) {
	memory := &MemoryFlowPersist{Sessions: map[any]*Flow{}}

	return memory, nil
}

func (memory *MemoryFlowPersist) Set(ID any, flow *Flow) error {
	memory.Sessions[ID] = flow

	return nil
}

func (memory *MemoryFlowPersist) Get(ID any) (*Flow, error) {
	flow := memory.Sessions[ID]

	if flow == nil {
		return nil, fmt.Errorf("not found flow")
	}

	return flow, nil
}
