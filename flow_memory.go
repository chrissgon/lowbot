package lowbot

type MemoryFlowPersist struct {
	Sessions map[any]*Flow
}

func NewMemoryFlowPersist() IFlowPersist {
	memory := &MemoryFlowPersist{Sessions: map[any]*Flow{}}

	return memory
}

func (memory *MemoryFlowPersist) Set(ID any, flow *Flow) error {
	memory.Sessions[ID] = flow

	return nil
}

func (memory *MemoryFlowPersist) Get(ID any) (*Flow, error) {
	flow := memory.Sessions[ID]

	if flow == nil {
		return nil, ERR_NIL_FLOW
	}

	return flow, nil
}
