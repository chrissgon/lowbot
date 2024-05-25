package lowbot

import "github.com/google/uuid"

type IConsumer interface {
	Run(*Interaction, IChannel) error
}

type Consumer struct {
	ConsumerID uuid.UUID
	Name       string
}
