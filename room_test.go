package lowbot

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

var ROOM_GUESTS_MOCK = RoomGuests{
	WHO_MOCK.WhoID: NewGuest(WHO_MOCK, CHANNEL_MOCK),
}
var ROOM_MOCK = NewRoom(ROOM_GUESTS_MOCK)

func TestRoom_NewRoom(t *testing.T) {
	guests := RoomGuests{
		WHO_MOCK.WhoID: NewGuest(WHO_MOCK, CHANNEL_MOCK),
	}

	expect := &Room{
		RoomID:       uuid.New(),
		Guests:       guests,
		Interactions: []*Interaction{},
	}
	have := NewRoom(guests)

	have.RoomID = expect.RoomID

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestRoom_AddGuest(t *testing.T) {
	room := NewRoom(ROOM_GUESTS_MOCK)

	have := len(room.Guests)
	expect := 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	room.AddGuest(NewGuest(NewWho("2", "amanda"), CHANNEL_MOCK))

	have = len(room.Guests)
	expect = 2

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestRoom_AddInteraction(t *testing.T) {
	channelCount = 0

	room := NewRoom(ROOM_GUESTS_MOCK)
	room.AddGuest(NewGuest(NewWho("2", "amanda"), CHANNEL_MOCK))
	room.AddGuest(NewGuest(NewWho("3", "thais"), CHANNEL_MOCK))

	have := len(room.Interactions)
	expect := 0

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	room.AddInteraction(NewInteractionMessageText(WHO_MOCK, WHO_MOCK, TEXT))

	have = len(room.Interactions)
	expect = 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	have = channelCount
	expect = 2

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestRoom_SendInteractionExcludingSender(t *testing.T) {
	channelCount = 0

	room := NewRoom(ROOM_GUESTS_MOCK)
	room.AddGuest(NewGuest(NewWho("2", "amanda"), CHANNEL_MOCK))
	room.AddGuest(NewGuest(NewWho("3", "thais"), CHANNEL_MOCK))

	interaction := NewInteractionMessageText(WHO_MOCK, WHO_MOCK, TEXT)

	errs := room.SendInteractionExcludingSender(interaction)

	if (len(errs) == 0){
		t.Errorf(FormatTestError(1, len(errs)))
	}

	have := channelCount
	expect := 2

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
