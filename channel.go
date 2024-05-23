package lowbot

import "github.com/google/uuid"

type IChannel interface {
	ChannelID() uuid.UUID
	SendAudio(*Interaction) error
	SendButton(*Interaction) error
	SendDocument(*Interaction) error
	SendImage(*Interaction) error
	SendText(*Interaction) error
	SendVideo(*Interaction) error
	Next(chan *Interaction)
}

type Interaction struct {
	SessionID  string
	Type       InteractionType
	Parameters InteractionParameters
	Custom     map[string]any
}

type InteractionParameters StepParameters

type InteractionType string

const (
	MESSAGE_AUDIO    InteractionType = "audio"
	MESSAGE_BUTTON   InteractionType = "button"
	MESSAGE_DOCUMENT InteractionType = "document"
	MESSAGE_IMAGE    InteractionType = "image"
	MESSAGE_TEXT     InteractionType = "text"
	MESSAGE_VIDEO    InteractionType = "video"
	EVENT_TYPING     InteractionType = "typing"
)

func NewInteractionMessageAudio(sessionID string, audio string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE_AUDIO,
		Parameters: InteractionParameters{
			Text:  text,
			Audio: audio,
		},
	}
}

func NewInteractionMessageButton(sessionID string, buttons []string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE_BUTTON,
		Parameters: InteractionParameters{
			Text:    text,
			Buttons: buttons,
		},
	}
}

func NewInteractionMessageDocument(sessionID string, document string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE_DOCUMENT,
		Parameters: InteractionParameters{
			Text:     text,
			Document: document,
		},
	}
}

func NewInteractionMessageImage(sessionID string, image string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE_IMAGE,
		Parameters: InteractionParameters{
			Text:  text,
			Image: image,
		},
	}
}

func NewInteractionMessageText(sessionID string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE_TEXT,
		Parameters: InteractionParameters{
			Text: text,
		},
	}
}

func NewInteractionMessageVideo(sessionID string, video string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE_VIDEO,
		Parameters: InteractionParameters{
			Text:  text,
			Video: video,
		},
	}
}
