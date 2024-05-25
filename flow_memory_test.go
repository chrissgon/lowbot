package lowbot

import (
	"reflect"
	"testing"
)

func TestFlowMemory_NewMemoryFlowPersist(t *testing.T) {
	expect := &MemoryFlowPersist{Sessions: map[any]*Flow{}}
	have, err := NewMemoryFlowPersist()

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFlowMemory_Set(t *testing.T) {
	ID := 1
	flow := newFlowMock()
	persist, _ := NewMemoryFlowPersist()

	err := persist.Set(ID, &flow)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}
}

func TestFlowMemory_Get(t *testing.T) {
	ID := 1
	flow := newFlowMock()
	persist, _ := NewMemoryFlowPersist()

	persist.Set(ID, &flow)

	expect := &flow
	have, err := persist.Get(ID)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}
