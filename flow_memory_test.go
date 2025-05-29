package lowbot

import (
	"errors"
	"reflect"
	"testing"
)

func TestFlowMemory_NewMemoryFlowPersist(t *testing.T) {
	expect := &MemoryFlowPersist{Sessions: map[any]*Flow{}}
	have := NewMemoryFlowPersist()

	if !reflect.DeepEqual(expect, have) {
		t.Error(FormatTestError(expect, have))
	}
}

func TestFlowMemory_Set(t *testing.T) {
	ID := 1
	flow := newFlowMock()
	persist := NewMemoryFlowPersist()

	err := persist.Set(ID, flow)

	if err != nil {
		t.Error(FormatTestError(nil, err))
	}
}

func TestFlowMemory_Get(t *testing.T) {
	ID := 1
	flow := newFlowMock()
	persist := NewMemoryFlowPersist()

	_, err := persist.Get(ID)

	if !errors.Is(err, ERR_NIL_FLOW) {
		t.Error(FormatTestError(ERR_NIL_FLOW, err))
	}

	persist.Set(ID, flow)

	expect := flow
	have, err := persist.Get(ID)

	if err != nil {
		t.Error(FormatTestError(nil, err))
	}

	if !reflect.DeepEqual(expect, have) {
		t.Error(FormatTestError(expect, have))
	}
}
