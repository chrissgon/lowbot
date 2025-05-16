package lowbot

import (
	"errors"
	"testing"
)

func TestActions_SetCustomActions(t *testing.T) {
	name := "Custom"

	SetCustomActions(ActionsMap{
		name: func(i *Interaction) (*Interaction, bool) { return nil, true },
	})

	_, exists := actions[name]

	if !exists {
		t.Fatalf("not found action: %s", name)
	}
}

func TestActions_GetAction(t *testing.T) {
	_, err := GetAction(nil)

	if !errors.Is(err, ERR_NIL_FLOW) {
		t.Fatal(FormatTestError(ERR_NIL_FLOW, err))
	}

	flow := newFlowMock()
	flow.CurrentStep = nil

	_, err = GetAction(flow)

	if !errors.Is(err, ERR_NIL_STEP) {
		t.Fatal(FormatTestError(ERR_NIL_STEP, err))
	}

	flow.Steps["audio"].Action = "undefined"
	flow.CurrentStep = flow.Steps["audio"]

	_, err = GetAction(flow)

	if !errors.Is(err, ERR_UNKNOWN_ACTION) {
		t.Fatal(FormatTestError(ERR_UNKNOWN_ACTION, err))
	}

	flow.CurrentStep = flow.Steps["button"]

	_, err = GetAction(flow)

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}
}

func TestActions_RunActionButton(t *testing.T) {
	interaction := NewInteractionMessageText("")
	flow := newFlowMock()

	flow.Steps["button"].Parameters.Buttons = []string{"yes", "no"}
	flow.Steps["button"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["button"]

	action, err := GetAction(flow)

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}

	interaction, wait := action(interaction)

	if interaction == nil {
		t.Fatal(FormatTestError("interaction", nil))
	}

	if !wait {
		t.Fatal(FormatTestError(true, false))
	}

	if interaction.Type != MESSAGE_BUTTON {
		t.Fatal(FormatTestError(MESSAGE_BUTTON, interaction.Type))
	}
}

func TestActions_RunActionFileAudio(t *testing.T) {
	interaction := NewInteractionMessageText("")
	flow := newFlowMock()

	flow.Steps["audio"].Parameters.Path = "./mocks/music.mp3"
	flow.Steps["audio"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["audio"]

	action, err := GetAction(flow)

	interaction.Step = *flow.CurrentStep

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}

	interaction, wait := action(interaction)

	if interaction == nil {
		t.Fatal(FormatTestError("interaction", nil))
	}

	if wait {
		t.Fatal(FormatTestError(false, wait))
	}

	if interaction.Type != MESSAGE_FILE {
		t.Fatal(FormatTestError(MESSAGE_FILE, interaction.Type))
	}

	if !interaction.Parameters.File.GetFile().IsAudio() {
		t.Fatal(FormatTestError(FILETYPE_AUDIO, interaction.Parameters.File.GetFile().FileType))
	}
}

func TestActions_RunActionFileDocument(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText("")
	flow := newFlowMock()

	flow.Steps["document"].Parameters.Path = "./mocks/features.txt"
	flow.Steps["document"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["document"]

	action, err := GetAction(flow)

	interaction.Step = *flow.CurrentStep

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}

	interaction, wait := action(interaction)

	if interaction == nil {
		t.Fatal(FormatTestError("interaction", nil))
	}

	if wait {
		t.Fatal(FormatTestError(false, wait))
	}

	if interaction.Type != MESSAGE_FILE {
		t.Fatal(FormatTestError(MESSAGE_FILE, interaction.Type))
	}

	if !interaction.Parameters.File.GetFile().IsDocument() {
		t.Fatal(FormatTestError(FILETYPE_DOCUMENT, interaction.Parameters.File.GetFile().FileType))
	}
}

func TestActions_RunActionFileImage(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText("")
	flow := newFlowMock()

	flow.Steps["image"].Parameters.Path = "./mocks/image.jpg"
	flow.Steps["image"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["image"]

	action, err := GetAction(flow)

	interaction.Step = *flow.CurrentStep

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}

	interaction, wait := action(interaction)

	if interaction == nil {
		t.Fatal(FormatTestError("interaction", nil))
	}

	if wait {
		t.Fatal(FormatTestError(false, wait))
	}

	if interaction.Type != MESSAGE_FILE {
		t.Fatal(FormatTestError(MESSAGE_FILE, interaction.Type))
	}

	if !interaction.Parameters.File.GetFile().IsImage() {
		t.Fatal(FormatTestError(FILETYPE_IMAGE, interaction.Parameters.File.GetFile().FileType))
	}
}

func TestActions_RunActionFileVideo(t *testing.T) {
	channelCount = 0
	interaction := NewInteractionMessageText("")
	flow := newFlowMock()

	flow.Steps["video"].Parameters.Path = "./mocks/video.mp4"
	flow.Steps["video"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["video"]

	action, err := GetAction(flow)

	interaction.Step = *flow.CurrentStep

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}

	interaction, wait := action(interaction)

	if interaction == nil {
		t.Fatal(FormatTestError("interaction", nil))
	}

	if wait {
		t.Fatal(FormatTestError(false, wait))
	}

	if interaction.Type != MESSAGE_FILE {
		t.Fatal(FormatTestError(MESSAGE_FILE, interaction.Type))
	}

	if !interaction.Parameters.File.GetFile().IsVideo() {
		t.Fatal(FormatTestError(FILETYPE_VIDEO, interaction.Parameters.File.GetFile().FileType))
	}
}

func TestActions_RunActionInput(t *testing.T) {
	channelTriggerError = false
	channelCount = 0
	interaction := NewInteractionMessageText("")
	flow := newFlowMock()

	flow.Steps["input"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["input"]

	action, err := GetAction(flow)

	interaction.Step = *flow.CurrentStep

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}

	interaction, wait := action(interaction)

	if interaction == nil {
		t.Fatal(FormatTestError("interaction", nil))
	}

	if !wait {
		t.Fatal(FormatTestError(true, wait))
	}

	if interaction.Type != MESSAGE_TEXT {
		t.Fatal(FormatTestError(MESSAGE_FILE, interaction.Type))
	}
}

func TestActions_RunActionText(t *testing.T) {
	channelTriggerError = false
	channelCount = 0
	interaction := NewInteractionMessageText("")
	flow := newFlowMock()

	flow.Steps["text"].Parameters.Texts = []string{"text"}
	flow.CurrentStep = flow.Steps["text"]

	action, err := GetAction(flow)

	interaction.Step = *flow.CurrentStep

	if err != nil {
		t.Fatal(FormatTestError(nil, err))
	}

	interaction, wait := action(interaction)

	if interaction == nil {
		t.Fatal(FormatTestError("interaction", nil))
	}

	if wait {
		t.Fatal(FormatTestError(false, wait))
	}

	if interaction.Type != MESSAGE_TEXT {
		t.Fatal(FormatTestError(MESSAGE_FILE, interaction.Type))
	}
}
