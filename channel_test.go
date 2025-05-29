package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

var (
	channelCount                           = 0
	channelLastMethodCalled                = ""
	channelLastInteractionSent Interaction = Interaction{}
	channelTriggerError                    = true

	ErrMock = fmt.Errorf("error mock")
)

type mockChannel struct {
	*Channel
	startedTimes int
	stoppedTimes int
	fail         bool
}

var CHANNEL_MOCK = newMockChannel()

func newMockChannel() IChannel {
	return &mockChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      "mock channel",
			Broadcast: NewBroadcast[Interaction](),
		},
		startedTimes: 0,
		stoppedTimes: 0,
		fail:         false,
	}
}

func (m *mockChannel) GetChannel() *Channel {
	return m.Channel
}

func (m *mockChannel) Stop() error {
	m.stoppedTimes++
	if m.fail {
		return ErrMock
	}
	return nil
}

func (m *mockChannel) Start() error {
	m.startedTimes++
	if m.fail {
		return ErrMock
	}
	return nil
}

func (m *mockChannel) Next() {
}

func (m *mockChannel) SendAudio(interaction Interaction) error {
	channelLastMethodCalled = "SendAudio"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendButton(interaction Interaction) error {
	channelLastMethodCalled = "SendButton"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendDocument(interaction Interaction) error {
	channelLastMethodCalled = "SendDocument"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendImage(interaction Interaction) error {
	channelLastMethodCalled = "SendImage"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

func (m *mockChannel) SendText(interaction Interaction) error {
	channelLastMethodCalled = "SendText"
	channelCount++
	channelLastInteractionSent = interaction
	if channelTriggerError {
		return ErrMock
	}
	return nil
}

func (m *mockChannel) SendVideo(interaction Interaction) error {
	channelLastMethodCalled = "SendVideo"
	channelCount++
	channelLastInteractionSent = interaction
	return nil
}

// func TestChannel_SendInteractionText(t *testing.T) {
// 	channelTriggerError = true
// 	channelCount = 0

// 	interaction := NewInteractionMessageText(DESTINATION_MOCK, SENDER_MOCK, TEXT)

// 	err := SendInteraction(CHANNEL_MOCK, interaction)

// 	if !errors.Is(err, ErrMock) {
// 		t.Error(FormatTestError(ErrMock, err))
// 	}

// 	if channelCount != 1 {
// 		t.Error(FormatTestError(1, channelCount))
// 	}

// 	if channelLastMethodCalled != "SendText" {
// 		t.Error(FormatTestError("SendText", channelLastMethodCalled))
// 	}
// }

// func TestChannel_SendInteractionButton(t *testing.T) {
// 	channelCount = 0

// 	interaction := NewInteractionMessageButton(DESTINATION_MOCK, SENDER_MOCK, BUTTONS, TEXT)

// 	err := SendInteraction(CHANNEL_MOCK, interaction)

// 	if err != nil {
// 		t.Error(FormatTestError(nil, err))
// 	}

// 	if channelCount != 1 {
// 		t.Error(FormatTestError(1, channelCount))
// 	}

// 	if channelLastMethodCalled != "SendButton" {
// 		t.Error(FormatTestError("SendButton", channelLastMethodCalled))
// 	}
// }

// func TestChannel_SendInteractionAudio(t *testing.T) {
// 	channelCount = 0

// 	interaction := NewInteractionMessageFile(DESTINATION_MOCK, SENDER_MOCK, "./mocks/audio.mp3", TEXT)

// 	err := SendInteraction(CHANNEL_MOCK, interaction)

// 	if err != nil {
// 		t.Error(FormatTestError(nil, err))
// 	}

// 	if channelCount != 1 {
// 		t.Error(FormatTestError(1, channelCount))
// 	}

// 	if channelLastMethodCalled != "SendAudio" {
// 		t.Error(FormatTestError("SendAudio", channelLastMethodCalled))
// 	}
// }

// func TestChannel_SendInteractionDocument(t *testing.T) {
// 	channelCount = 0

// 	interaction := NewInteractionMessageFile(DESTINATION_MOCK, SENDER_MOCK, "./mocks/features.txt", TEXT)

// 	err := SendInteraction(CHANNEL_MOCK, interaction)

// 	if err != nil {
// 		t.Error(FormatTestError(nil, err))
// 	}

// 	if channelCount != 1 {
// 		t.Error(FormatTestError(1, channelCount))
// 	}

// 	if channelLastMethodCalled != "SendDocument" {
// 		t.Error(FormatTestError("SendDocument", channelLastMethodCalled))
// 	}
// }

// func TestChannel_SendInteractionImage(t *testing.T) {
// 	channelCount = 0

// 	interaction := NewInteractionMessageFile(DESTINATION_MOCK, SENDER_MOCK, "./mocks/image.jpg", TEXT)

// 	err := SendInteraction(CHANNEL_MOCK, interaction)

// 	if err != nil {
// 		t.Error(FormatTestError(nil, err))
// 	}

// 	if channelCount != 1 {
// 		t.Error(FormatTestError(1, channelCount))
// 	}

// 	if channelLastMethodCalled != "SendImage" {
// 		t.Error(FormatTestError("SendImage", channelLastMethodCalled))
// 	}
// }

// func TestChannel_SendInteractionVideo(t *testing.T) {
// 	channelCount = 0

// 	interaction := NewInteractionMessageFile(DESTINATION_MOCK, SENDER_MOCK, "./mocks/video.mp4", TEXT)

// 	err := SendInteraction(CHANNEL_MOCK, interaction)

// 	if err != nil {
// 		t.Error(FormatTestError(nil, err))
// 	}

// 	if channelCount != 1 {
// 		t.Error(FormatTestError(1, channelCount))
// 	}

// 	if channelLastMethodCalled != "SendVideo" {
// 		t.Error(FormatTestError("SendVideo", channelLastMethodCalled))
// 	}
// }
