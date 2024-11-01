package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

type ConsumerChannel struct {
	*Channel
	Consumer IConsumer
	running  bool
}

func NewConsumerChannel(consumer IConsumer) (IChannel, error) {
	return &ConsumerChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      consumer.GetConsumer().Name,
			Broadcast: NewBroadcast[*Interaction](),
		},
		Consumer: consumer,
		running:  false,
	}, nil
}

func (channel *ConsumerChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *ConsumerChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	channel.running = true

	return nil
}

func (channel *ConsumerChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	channel.running = false

	return nil
}

func (channel *ConsumerChannel) SendAudio(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

func (channel *ConsumerChannel) SendButton(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

func (channel *ConsumerChannel) SendDocument(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

func (channel *ConsumerChannel) SendImage(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

func (channel *ConsumerChannel) SendText(interaction *Interaction) error {
	answersInteraction, err := channel.Consumer.Run(interaction)

	
	if err != nil {
		return err
	}
	
	for _, answerInteraction := range answersInteraction {
		fmt.Println("sendText", answerInteraction.Sender.Name)
		channel.Broadcast.Send(answerInteraction)
	}

	return err
	// fmt.Println("sendText", interaction.Parameters.Text)
	// _, err := channel.Consumer.Run(interaction)
	// return err
}

func (channel *ConsumerChannel) SendVideo(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}
