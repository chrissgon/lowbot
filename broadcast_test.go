package lowbot

import (
	"reflect"
	"testing"
	"time"
)

func TestBroadcast_NewBroadcast(t *testing.T) {
	have := &Broadcast[int]{
		listeners: []chan int{},
	}
	expect := NewBroadcast[int]()

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestBroadcast_Listen(t *testing.T) {
	broadcast := NewBroadcast[int]()

	if len(broadcast.listeners) != 0 {
		t.Errorf(FormatTestError(0, len(broadcast.listeners)))
	}

	channel := broadcast.Listen()

	if len(broadcast.listeners) != 1 {
		t.Errorf(FormatTestError(1, len(broadcast.listeners)))
	}
	if channel == nil {
		t.Errorf(FormatTestError("not nil channel", nil))
	}
}

func TestBroadcast_Send(t *testing.T) {
	passed := false
	broadcast := NewBroadcast[bool]()

	channel := broadcast.Listen()

	go func() {
		for {
			passed = <-channel
		}
	}()

	broadcast.Send(true)
	time.Sleep(1 * time.Millisecond)
	broadcast.Send(true)
	time.Sleep(1 * time.Millisecond)

	if !passed {
		t.Errorf(FormatTestError(true, false))
	}
}
func TestBroadcast_Close(t *testing.T) {
	broadcast := NewBroadcast[bool]()

	broadcast.Listen()

	broadcast.Close()

	if len(broadcast.listeners) != 0 {
		t.Errorf(FormatTestError(0, len(broadcast.listeners)))
	}
}
