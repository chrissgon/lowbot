package lowbot

import (
	"reflect"
	"testing"
)

var WHO_MOCK = NewWho("1", "chris")

func TestWho_NewWho(t *testing.T) {
	expect := &Who{
		WhoID:  "1",
		Name:   "chris",
		Custom: map[string]any{},
	}
	have := NewWho("1", "chris")

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}
