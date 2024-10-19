package lowbot

import (
	"github.com/google/uuid"
)

type IConsumer interface {
	Run(*Interaction) ([]*Interaction, error)
	GetConsumer() *Consumer
}

type Consumer struct {
	ConsumerID uuid.UUID
	Name       string
}
