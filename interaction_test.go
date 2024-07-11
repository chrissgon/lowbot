package lowbot

import (
	"path/filepath"
	"reflect"
	"testing"
)

const (
	BUTTON = "button"
	TEXT   = "text"
	FILE   = "file"
)

var (
	BUTTONS          = []string{"button"}
	DESTINATION_MOCK = NewWho("1", "chris")
	SENDER_MOCK      = NewWho("2", "amanda")
)

func TestInteraction_NewInteractionMessageButton(t *testing.T) {
	have := NewInteractionMessageButton(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, BUTTONS, TEXT)

	if CHANNEL_MOCK.GetChannel() != have.Channel {
		t.Errorf(FormatTestError(CHANNEL_MOCK.GetChannel(), have.Channel))
	}

	if DESTINATION_MOCK != have.Destination {
		t.Errorf(FormatTestError(DESTINATION_MOCK, have.Destination))
	}

	if SENDER_MOCK != have.Sender {
		t.Errorf(FormatTestError(SENDER_MOCK, have.Sender))
	}

	if MESSAGE_BUTTON != have.Type {
		t.Errorf(FormatTestError(MESSAGE_BUTTON, have.Type))
	}

	if TEXT != have.Parameters.Text {
		t.Errorf(FormatTestError(TEXT, have.Parameters.Text))
	}

	if len(BUTTONS) != len(have.Parameters.Buttons) {
		t.Errorf(FormatTestError(BUTTONS, have.Parameters.Buttons))
	}
}

func TestInteraction_NewInteractionMessageFile(t *testing.T) {
	have := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, FILE, TEXT)

	if CHANNEL_MOCK.GetChannel() != have.Channel {
		t.Errorf(FormatTestError(CHANNEL_MOCK.GetChannel(), have.Channel))
	}

	if DESTINATION_MOCK != have.Destination {
		t.Errorf(FormatTestError(DESTINATION_MOCK, have.Destination))
	}

	if SENDER_MOCK != have.Sender {
		t.Errorf(FormatTestError(SENDER_MOCK, have.Sender))
	}

	if MESSAGE_FILE != have.Type {
		t.Errorf(FormatTestError(MESSAGE_FILE, have.Type))
	}

	if TEXT != have.Parameters.Text {
		t.Errorf(FormatTestError(TEXT, have.Parameters.Text))
	}

	abs, _ := filepath.Abs(FILE)

	if abs != have.Parameters.File.GetFile().Path {
		t.Errorf(FormatTestError(abs, have.Parameters.File.GetFile().Path))
	}
}

func TestInteraction_NewInteractionMessageText(t *testing.T) {
	have := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, TEXT)

	if CHANNEL_MOCK.GetChannel() != have.Channel {
		t.Errorf(FormatTestError(CHANNEL_MOCK.GetChannel(), have.Channel))
	}

	if DESTINATION_MOCK != have.Destination {
		t.Errorf(FormatTestError(DESTINATION_MOCK, have.Destination))
	}

	if SENDER_MOCK != have.Sender {
		t.Errorf(FormatTestError(SENDER_MOCK, have.Sender))
	}

	if MESSAGE_TEXT != have.Type {
		t.Errorf(FormatTestError(MESSAGE_TEXT, have.Type))
	}

	if TEXT != have.Parameters.Text {
		t.Errorf(FormatTestError(TEXT, have.Parameters.Text))
	}
}

func TestInteraction_SetReplier(t *testing.T) {
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, TEXT)

	interaction.SetReplier(WHO_MOCK)

	have := interaction.Replier
	expect := WHO_MOCK

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestInteraction_SetDestination(t *testing.T) {
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, TEXT)

	interaction.SetDestination(WHO_MOCK)

	have := interaction.Destination
	expect := WHO_MOCK

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}
