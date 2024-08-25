package lowbot

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

var (
	channelCount            = 0
	channelLastMethodCalled = ""
	channelLastInteractionSent *Interaction = nil
	channelTriggerError = true

	ERR_MOCK = errors.New("error mock")
)

type mockChannel struct {
	*Channel
}

var CHANNEL_MOCK = newMockChannel()

func newMockChannel() IChannel {
	return &mockChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      "mock channel",
			Broadcast: NewBroadcast[*Interaction](),
		},
	}
}

func (m *mockChannel) GetChannel() *Channel {
	return m.Channel
}

func (m *mockChannel) Close() error {
	return nil
}

func (m *mockChannel) Next() {
}

func (m *mockChannel) SendAudio(interaction *Interaction) error {
	channelLastMethodCalled = "SendAudio"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendButton(interaction *Interaction) error {
	channelLastMethodCalled = "SendButton"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendDocument(interaction *Interaction) error {
	channelLastMethodCalled = "SendDocument"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendImage(interaction *Interaction) error {
	channelLastMethodCalled = "SendImage"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendText(interaction *Interaction) error {
	channelLastMethodCalled = "SendText"
	channelCount++
	channelLastInteractionSent = interaction
	if(channelTriggerError){
		return ERR_MOCK
	}
	return nil
}

func (m *mockChannel) SendVideo(interaction *Interaction) error {
	channelLastMethodCalled = "SendVideo"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func TestChannel_SendInteractionText(t *testing.T) {
	channelTriggerError = true
	channelCount = 0

	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if !errors.Is(err, ERR_MOCK) {
		t.Errorf(FormatTestError(ERR_MOCK, err))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendText" {
		t.Errorf(FormatTestError("SendText", channelLastMethodCalled))
	}
}

func TestChannel_SendInteractionButton(t *testing.T) {
	channelCount = 0

	interaction := NewInteractionMessageButton(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, BUTTONS, TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendButton" {
		t.Errorf(FormatTestError("SendButton", channelLastMethodCalled))
	}
}

func TestChannel_SendInteractionAudio(t *testing.T) {
	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/audio.mp3", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendAudio" {
		t.Errorf(FormatTestError("SendAudio", channelLastMethodCalled))
	}
}

func TestChannel_SendInteractionDocument(t *testing.T) {
	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/features.txt", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendDocument" {
		t.Errorf(FormatTestError("SendDocument", channelLastMethodCalled))
	}
}

func TestChannel_SendInteractionImage(t *testing.T) {
	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/image.jpg", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendImage" {
		t.Errorf(FormatTestError("SendImage", channelLastMethodCalled))
	}
}

func TestChannel_SendInteractionVideo(t *testing.T) {
	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/video.mp4", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendVideo" {
		t.Errorf(FormatTestError("SendVideo", channelLastMethodCalled))
	}
}
