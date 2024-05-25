package lowbot

import (
	"github.com/google/uuid"
)

type IChannel interface {
	GetChannel() *Channel
	Next(chan *Interaction)
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
}
