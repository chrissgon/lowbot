package lowbot

import (
	"github.com/google/uuid"
)

type IChannel interface {
	ChannelID() uuid.UUID
	Next(chan *Interaction)
	SendAudio(*Interaction) error
	SendButton(*Interaction) error
	SendDocument(*Interaction) error
	SendImage(*Interaction) error
	SendText(*Interaction) error
	SendVideo(*Interaction) error
}

type Interaction struct {
	ChannelID  uuid.UUID
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

func SendInteraction(channel IChannel, interaction *Interaction) error {
	if interaction.Type == MESSAGE_AUDIO {
		return channel.SendAudio(interaction)
	}
	if interaction.Type == MESSAGE_BUTTON {
		return channel.SendButton(interaction)
	}
	if interaction.Type == MESSAGE_DOCUMENT {
		return channel.SendDocument(interaction)
	}
	if interaction.Type == MESSAGE_IMAGE {
		return channel.SendImage(interaction)
	}
	if interaction.Type == MESSAGE_TEXT {
		return channel.SendText(interaction)
	}
	if interaction.Type == MESSAGE_VIDEO {
		return channel.SendVideo(interaction)
	}

	return ERR_UNKNOWN_INTERACTION_TYPE
}

func NewInteractionMessageAudio(channelID uuid.UUID, sessionID string, audio string, text string) *Interaction {
	return &Interaction{
		ChannelID: channelID,
		SessionID: sessionID,
		Type:      MESSAGE_AUDIO,
		Parameters: InteractionParameters{
			Text:  text,
			Audio: audio,
		},
	}
}

func NewInteractionMessageButton(channelID uuid.UUID, sessionID string, buttons []string, text string) *Interaction {
	return &Interaction{
		ChannelID: channelID,
		SessionID: sessionID,
		Type:      MESSAGE_BUTTON,
		Parameters: InteractionParameters{
			Text:    text,
			Buttons: buttons,
		},
	}
}

func NewInteractionMessageDocument(channelID uuid.UUID, sessionID string, document string, text string) *Interaction {
	return &Interaction{
		ChannelID: channelID,
		SessionID: sessionID,
		Type:      MESSAGE_DOCUMENT,
		Parameters: InteractionParameters{
			Text:     text,
			Document: document,
		},
	}
}

func NewInteractionMessageImage(channelID uuid.UUID, sessionID string, image string, text string) *Interaction {
	return &Interaction{
		ChannelID: channelID,
		SessionID: sessionID,
		Type:      MESSAGE_IMAGE,
		Parameters: InteractionParameters{
			Text:  text,
			Image: image,
		},
	}
}

func NewInteractionMessageText(channelID uuid.UUID, sessionID string, text string) *Interaction {
	return &Interaction{
		ChannelID: channelID,
		SessionID: sessionID,
		Type:      MESSAGE_TEXT,
		Parameters: InteractionParameters{
			Text: text,
		},
	}
}

func NewInteractionMessageVideo(channelID uuid.UUID, sessionID string, video string, text string) *Interaction {
	return &Interaction{
		ChannelID: channelID,
		SessionID: sessionID,
		Type:      MESSAGE_VIDEO,
		Parameters: InteractionParameters{
			Text:  text,
			Video: video,
		},
	}
}
