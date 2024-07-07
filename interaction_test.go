package lowbot

import (
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

const (
	BUTTON = "button"
	TEXT   = "text"
	FILE   = "file"
)

var (
	BUTTONS      = []string{"button"}
	DESTINATION_MOCK     = NewWho(uuid.New().String(), "WHO DESTINATION")
	SENDER_MOCK     = NewWho(uuid.New().String(), "WHO SENDER")
	CHANNEL_MOCK = newMockChannel()
)

type mockChannel struct {
	*Channel
}

func newMockChannel() IChannel {
	return &mockChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      "mock",
		},
	}
}

func (m *mockChannel) GetChannel() *Channel {
	return m.Channel
}

func (m *mockChannel) Close() error {
	return nil
}

func (m *mockChannel) Next(interaction chan *Interaction) {
	interaction <- NewInteractionMessageText(m, DESTINATION_MOCK, SENDER_MOCK, TEXT)
}

func (m *mockChannel) SendAudio(*Interaction) error {
	return nil
}

func (m *mockChannel) SendButton(*Interaction) error {
	return nil
}

func (m *mockChannel) SendDocument(*Interaction) error {
	return nil
}

func (m *mockChannel) SendImage(*Interaction) error {
	return nil
}

func (m *mockChannel) SendText(*Interaction) error {
	return nil
}

func (m *mockChannel) SendVideo(*Interaction) error {
	return nil
}

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
