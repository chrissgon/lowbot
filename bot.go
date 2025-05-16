package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

type Bot struct {
	BotID    uuid.UUID
	Consumer IConsumer
	Channels map[uuid.UUID]IChannel
	Running  bool
}

func NewBot(consumer IConsumer, channels map[uuid.UUID]IChannel) *Bot {
	return &Bot{
		BotID:    uuid.New(),
		Consumer: consumer,
		Channels: channels,
		Running:  false,
	}
}

func (bot *Bot) Start() error {
	for _, channel := range bot.Channels {
		err := bot.StartChannel(channel)

		if err != nil {
			return err
		}

		go bot.StartConsumerChannel(channel)
	}

	bot.Running = true

	return nil
}

func (bot *Bot) StartChannel(channel IChannel) error {
	return channel.Start()
}

func (bot *Bot) StartConsumerChannel(channel IChannel) {
	listener := channel.GetChannel().Broadcast.Listen()

	for interaction := range listener {
		answersInteraction, err := bot.Consumer.Run(interaction)

		if answersInteraction == nil {
			continue
		}

		for _, answerInteraction := range answersInteraction {
			SendInteraction(channel, answerInteraction)
		}

		// TODO: improve how to receive the consumer errors
		if err != nil {
			PrintLog(fmt.Sprintf("%v: WhoID:<%v> ERR: %v\n", bot.Consumer.GetConsumer().Name, interaction.From.WhoID, err))
		}
	}
}

func (bot *Bot) Stop() error {
	for _, channel := range bot.Channels {
		err := channel.Stop()

		if err != nil {
			return err
		}
	}

	bot.Running = false

	return nil
}
