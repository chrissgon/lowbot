package lowbot

import (
	"errors"
	"testing"
)

func TestActions_SetCustomActions(t *testing.T) {
	name := "Custom"

	SetCustomActions(ActionsMap{
		name: func(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) { return true, nil },
	})

	_, exists := actions[name]

	if !exists {
		t.Errorf("not found action: %s", name)
	}
}

func TestActions_RunAction(t *testing.T) {
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")

	_, err := RunAction(nil, interaction, nil)

	if !errors.Is(err, ERR_NIL_FLOW) {
		t.Errorf(FormatTestError(ERR_NIL_FLOW, err))
	}

	flow := newFlowMock()
	flow.CurrentStep = nil

	_, err = RunAction(flow, interaction, nil)

	if !errors.Is(err, ERR_NIL_STEP) {
		t.Errorf(FormatTestError(ERR_NIL_STEP, err))
	}

	flow.Steps["audio"].Action = "undefined"
	flow.CurrentStep = flow.Steps["audio"]

	_, err = RunAction(flow, interaction, nil)

	if !errors.Is(err, ERR_UNKNOWN_ACTION) {
		t.Errorf(FormatTestError(ERR_UNKNOWN_ACTION, err))
	}

	flow.CurrentStep = flow.Steps["button"]
	_, err = RunAction(flow, interaction, nil)

	if !errors.Is(err, ERR_NIL_CHANNEL) {
		t.Errorf(FormatTestError(ERR_NIL_CHANNEL, err))
	}

	_, err = RunAction(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}
}

func TestActions_RunActionButton(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["button"].Parameters.Buttons = []string{"yes", "no"}
	flow.Steps["button"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["button"]

	wait, err := RunActionButton(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if !wait {
		t.Errorf(FormatTestError(true, false))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendButton" {
		t.Errorf(FormatTestError("SendButton", channelLastMethodCalled))
	}

	if len(channelLastInteractionSent.Parameters.Buttons) != 2 {
		t.Errorf(FormatTestError(2, len(channelLastInteractionSent.Parameters.Buttons)))
	}

	if channelLastInteractionSent.Parameters.Buttons[0] != "yes" {
		t.Errorf(FormatTestError("yes", channelLastInteractionSent.Parameters.Buttons[0]))
	}

	if channelLastInteractionSent.Parameters.Buttons[1] != "no" {
		t.Errorf(FormatTestError("no", channelLastInteractionSent.Parameters.Buttons[1]))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}
}

func TestActions_RunActionFileAudio(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["audio"].Parameters.Path = "./mocks/music.mp3"
	flow.Steps["audio"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["audio"]

	wait, err := RunActionFile(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if wait {
		t.Errorf(FormatTestError(false, wait))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendAudio" {
		t.Errorf(FormatTestError("SendAudio", channelLastMethodCalled))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}

	if channelLastInteractionSent.Parameters.File == nil {
		t.Errorf(FormatTestError(File{}, nil))
	}
}

func TestActions_RunActionFileDocument(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["document"].Parameters.Path = "./mocks/features.txt"
	flow.Steps["document"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["document"]

	wait, err := RunActionFile(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if wait {
		t.Errorf(FormatTestError(false, wait))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendDocument" {
		t.Errorf(FormatTestError("SendDocument", channelLastMethodCalled))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}

	if channelLastInteractionSent.Parameters.File == nil {
		t.Errorf(FormatTestError(File{}, nil))
	}
}

func TestActions_RunActionFileImage(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["image"].Parameters.Path = "./mocks/image.jpg"
	flow.Steps["image"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["image"]

	wait, err := RunActionFile(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if wait {
		t.Errorf(FormatTestError(false, wait))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendImage" {
		t.Errorf(FormatTestError("SendImage", channelLastMethodCalled))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}

	if channelLastInteractionSent.Parameters.File == nil {
		t.Errorf(FormatTestError(File{}, nil))
	}
}

func TestActions_RunActionFileVideo(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["video"].Parameters.Path = "./mocks/video.mp4"
	flow.Steps["video"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["video"]

	wait, err := RunActionFile(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if wait {
		t.Errorf(FormatTestError(false, wait))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendVideo" {
		t.Errorf(FormatTestError("SendVideo", channelLastMethodCalled))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}

	if channelLastInteractionSent.Parameters.File == nil {
		t.Errorf(FormatTestError(File{}, nil))
	}
}

func TestActions_RunActionInput(t *testing.T) {
	channelTriggerError = false
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["button"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["button"]

	wait, err := RunActionInput(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if !wait {
		t.Errorf(FormatTestError(true, false))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendText" {
		t.Errorf(FormatTestError("SendText", channelLastMethodCalled))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}
}

func TestActions_RunActionText(t *testing.T) {
	channelTriggerError = false
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["button"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["button"]

	wait, err := RunActionText(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if wait {
		t.Errorf(FormatTestError(false, true))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendText" {
		t.Errorf(FormatTestError("SendText", channelLastMethodCalled))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}
}

func TestActions_RunActionWait(t *testing.T) {
	channelTriggerError = false
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["button"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["button"]

	wait, err := RunActionWait(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if !wait {
		t.Errorf(FormatTestError(true, false))
	}
}

func TestActions_RunActionRoom(t *testing.T) {
	channelTriggerError = false
	channelCount = 0
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow := newFlowMock()

	flow.Steps["button"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["button"]

	wait, err := RunActionRoom(flow, interaction, CHANNEL_MOCK)

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if wait {
		t.Errorf(FormatTestError(false, true))
	}

	if channelCount != 1 {
		t.Errorf(FormatTestError(1, channelCount))
	}

	if channelLastMethodCalled != "SendText" {
		t.Errorf(FormatTestError("SendText", channelLastMethodCalled))
	}

	if channelLastInteractionSent.Parameters.Text != "text" {
		t.Errorf(FormatTestError("text", channelLastInteractionSent.Parameters.Text))
	}

	if interaction.Custom["RoomID"] == nil {
		t.Errorf(FormatTestError("uuid", interaction.Custom["RoomID"]))
	}

	if len(roomManager.Rooms) != 1 {
		t.Errorf(FormatTestError(1, len(roomManager.Rooms)))
	}
}
