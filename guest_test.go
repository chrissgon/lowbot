package lowbot

import (
	"reflect"
	"testing"
)

func TestGuest_NewGuest(t *testing.T) {
	expect := &Guest{
		Who:     WHO_MOCK,
		Channel: CHANNEL_MOCK,
	}
	have := NewGuest(WHO_MOCK, CHANNEL_MOCK)

	if !reflect.DeepEqual(expect, have) {
		t.Error(FormatTestError(expect, have))
	}
}
