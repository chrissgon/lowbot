package lowbot

import "github.com/google/uuid"

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

// GetChannel implements IChannel.
func (channel *ConsumerChannel) GetChannel() *Channel {
	return channel.Channel
}

// Start implements IChannel.
func (channel *ConsumerChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	channel.running = true

	return nil
}

// Stop implements IChannel.
func (channel *ConsumerChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	channel.running = false

	return nil
}

// SendAudio implements IChannel.
func (channel *ConsumerChannel) SendAudio(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

// SendButton implements IChannel.
func (channel *ConsumerChannel) SendButton(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

// SendDocument implements IChannel.
func (channel *ConsumerChannel) SendDocument(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

// SendImage implements IChannel.
func (channel *ConsumerChannel) SendImage(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

// SendText implements IChannel.
func (channel *ConsumerChannel) SendText(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}

// SendVideo implements IChannel.
func (channel *ConsumerChannel) SendVideo(interaction *Interaction) error {
	_, err := channel.Consumer.Run(interaction)
	return err
}
