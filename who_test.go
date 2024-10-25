package lowbot

import (
	"reflect"
	"testing"
)

var WHO_MOCK = NewWho("1", "chris")

func TestWho_NewWho(t *testing.T) {
	have := NewWho("1", "chris")
	expect := &Who{
		WhoID:  "1",
		Name:   "chris",
		Custom: map[string]any{},
	}

	if !reflect.DeepEqual(expect, have) {
		t.Error(FormatTestError(expect, have))
	}
}
