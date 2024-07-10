package lowbot

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

var (
	channelCount            = 0
	channelLastMethodCalled = ""

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

func (m *mockChannel) SendAudio(*Interaction) error {
	channelLastMethodCalled = "SendAudio"
	channelCount++
	return nil
}

func (m *mockChannel) SendButton(*Interaction) error {
	channelLastMethodCalled = "SendButton"
	channelCount++
	return nil
}

func (m *mockChannel) SendDocument(*Interaction) error {
	channelLastMethodCalled = "SendDocument"
	channelCount++
	return nil
}

func (m *mockChannel) SendImage(*Interaction) error {
	channelLastMethodCalled = "SendImage"
	channelCount++
	return nil
}

func (m *mockChannel) SendText(*Interaction) error {
	channelLastMethodCalled = "SendText"
	channelCount++
	return ERR_MOCK
}

func (m *mockChannel) SendVideo(*Interaction) error {
	channelLastMethodCalled = "SendVideo"
	channelCount++
	return nil
}

func TestChannel_SendInteractionText(t *testing.T) {
	var have any
	var expect any

	channelCount = 0

	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if !errors.Is(err, ERR_MOCK) {
		t.Errorf(FormatTestError(ERR_MOCK, err))
	}

	have = channelCount
	expect = 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	have = channelLastMethodCalled
	expect = "SendText"

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestChannel_SendInteractionButton(t *testing.T) {
	var have any
	var expect any

	channelCount = 0

	interaction := NewInteractionMessageButton(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, BUTTONS, TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	have = channelCount
	expect = 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	have = channelLastMethodCalled
	expect = "SendButton"

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestChannel_SendInteractionAudio(t *testing.T) {
	var have any
	var expect any

	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/audio.mp3", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	have = channelCount
	expect = 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	have = channelLastMethodCalled
	expect = "SendAudio"

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestChannel_SendInteractionDocument(t *testing.T) {
	var have any
	var expect any

	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/features.txt", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	have = channelCount
	expect = 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	have = channelLastMethodCalled
	expect = "SendDocument"

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestChannel_SendInteractionImage(t *testing.T) {
	var have any
	var expect any

	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/image.jpg", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	have = channelCount
	expect = 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	have = channelLastMethodCalled
	expect = "SendImage"

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestChannel_SendInteractionVideo(t *testing.T) {
	var have any
	var expect any

	channelCount = 0

	interaction := NewInteractionMessageFile(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "./mocks/video.mp4", TEXT)

	err := SendInteraction(CHANNEL_MOCK, interaction)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	have = channelCount
	expect = 1

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

	have = channelLastMethodCalled
	expect = "SendVideo"

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
