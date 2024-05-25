package lowbot

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

var (
	CHANNELID = uuid.New()
)

const (
	SESSIONID = "41fd1dcb-3246-4886-883a-3a19db575687"
	AUDIO     = "audio"
	DOCUMENT  = "document"
	IMAGE     = "image"
	TEXT      = "text"
	VIDEO     = "video"
)

var BUTTONS = []string{"button"}

func TestNewInteractionMessageAudio(t *testing.T) {
	expect := Interaction{
		channelID: CHANNELID,
		SessionID: SESSIONID,
		Type:      MESSAGE_AUDIO,
		Parameters: InteractionParameters{
			Text:  TEXT,
			Audio: AUDIO,
		},
	}
	have := *NewInteractionMessageAudio(CHANNELID, SESSIONID, AUDIO, TEXT)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestNewInteractionMessageButton(t *testing.T) {
	expect := Interaction{
		channelID: CHANNELID,
		SessionID: SESSIONID,
		Type:      MESSAGE_BUTTON,
		Parameters: InteractionParameters{
			Text:    TEXT,
			Buttons: BUTTONS,
		},
	}
	have := *NewInteractionMessageButton(CHANNELID, SESSIONID, BUTTONS, TEXT)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestNewInteractionMessageDocument(t *testing.T) {
	expect := Interaction{
		channelID: CHANNELID,
		SessionID: SESSIONID,
		Type:      MESSAGE_DOCUMENT,
		Parameters: InteractionParameters{
			Text:     TEXT,
			Document: DOCUMENT,
		},
	}
	have := *NewInteractionMessageDocument(CHANNELID, SESSIONID, DOCUMENT, TEXT)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestNewInteractionMessageImage(t *testing.T) {
	expect := Interaction{
		channelID: CHANNELID,
		SessionID: SESSIONID,
		Type:      MESSAGE_IMAGE,
		Parameters: InteractionParameters{
			Text:  TEXT,
			Image: IMAGE,
		},
	}
	have := *NewInteractionMessageImage(CHANNELID, SESSIONID, IMAGE, TEXT)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestNewInteractionMessageText(t *testing.T) {
	expect := Interaction{
		channelID: CHANNELID,
		SessionID: SESSIONID,
		Type:      MESSAGE_TEXT,
		Parameters: InteractionParameters{
			Text: TEXT,
		},
	}
	have := *NewInteractionMessageText(CHANNELID, SESSIONID, TEXT)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestNewInteractionMessageVideo(t *testing.T) {
	expect := Interaction{
		channelID: CHANNELID,
		SessionID: SESSIONID,
		Type:      MESSAGE_VIDEO,
		Parameters: InteractionParameters{
			Text:  TEXT,
			Video: VIDEO,
		},
	}
	have := *NewInteractionMessageVideo(CHANNELID, SESSIONID, VIDEO, TEXT)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}
