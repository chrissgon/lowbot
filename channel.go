package lowbot

import (
	"github.com/google/uuid"
)

type IChannel interface {
	GetChannel() *Channel
	Stop() error
	Start() error
	SendAudio(*Interaction) error
	SendButton(*Interaction) error
	SendDocument(*Interaction) error
	SendImage(*Interaction) error
	SendText(*Interaction) error
	SendVideo(*Interaction) error
}

type Channel struct {
	ChannelID uuid.UUID
	Name      string
	Running   bool
	Broadcast *Broadcast[*Interaction]
}

func SendInteraction(channel IChannel, interaction *Interaction) error {
	if interaction.Type == MESSAGE_TEXT {
		return channel.SendText(interaction)
	}
	if interaction.Type == MESSAGE_BUTTON {
		return channel.SendButton(interaction)
	}

	if interaction.Parameters.File.IsAudio() {
		return channel.SendAudio(interaction)
	}
	if interaction.Parameters.File.IsImage() {
		return channel.SendImage(interaction)
	}
	if interaction.Parameters.File.IsVideo() {
		return channel.SendVideo(interaction)
	}

	return channel.SendDocument(interaction)
}
