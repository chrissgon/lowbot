package lowbot

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestBot_NewBot(t *testing.T) {
	bot, consumer, channel := newBotMock()

	if bot.Consumer == nil {
		t.Fatal("bot consumer should not be nil")
	}
	if bot.Consumer.GetConsumer().ConsumerID != consumer.GetConsumer().ConsumerID {
		t.Fatal("bot ConsumerID does not match with the one provided")
	}
	if bot.Channels == nil {
		t.Fatal("bot channels should not be nil")
	}
	if len(bot.Channels) != 1 {
		t.Fatalf("bot should has 1 channel, but got %v", len(bot.Channels))
	}

	_, channelExists := bot.Channels[channel.GetChannel().ChannelID]

	if !channelExists {
		t.Fatalf("bot should has the channel provided")
	}

	if bot.Running {
		t.Fatalf("bot should not be running yet")
	}
}

func TestBot_Start(t *testing.T) {
	bot, _, channel := newBotMock()

	bot.Start()

	time.Sleep(1 * time.Millisecond)

	if channel.(*mockChannel).startedTimes != 1 {
		t.Fatalf("bot should start channel once")
	}
	if len(channel.GetChannel().Broadcast.listeners) != 1 {
		t.Fatalf("bot should consumer channel once")
	}
	if !bot.Running {
		t.Fatalf("bot should be running")
	}

	channel.(*mockChannel).fail = true

	err := bot.Start()

	if !errors.Is(err, ErrMock) {
		t.Fatalf("bot should return an error when start")
	}
}

func TestBot_StartChannel(t *testing.T) {
	bot, _, channel := newBotMock()

	bot.StartChannel(channel)

	if channel.(*mockChannel).startedTimes != 1 {
		t.Fatalf("bot should start channel once")
	}
}

func TestBot_StartConsumerChannel(t *testing.T) {
	bot, consumer, channel := newBotMock()

	go bot.StartConsumerChannel(channel)

	time.Sleep(1 * time.Millisecond)

	if len(channel.GetChannel().Broadcast.listeners) != 1 {
		t.Fatalf("bot should consumer channel once")
	}

	interaction := NewInteractionMessageButton(DESTINATION_MOCK, SENDER_MOCK, BUTTONS, TEXT)
	channel.GetChannel().Broadcast.Send(interaction)

	time.Sleep(1 * time.Millisecond)

	if consumer.(*mockConsumer).ranTimes != 1 {
		t.Fatalf("bot should consumer an interaction")
	}
}

func TestBot_Stop(t *testing.T) {
	bot, _, channel := newBotMock()

	bot.Stop()

	if channel.(*mockChannel).stoppedTimes != 1 {
		t.Fatalf("bot should stop channel once")
	}
	if bot.Running {
		t.Fatalf("bot should not be running")
	}

	channel.(*mockChannel).fail = true

	err := bot.Stop()

	if !errors.Is(err, ErrMock) {
		t.Fatalf("bot should return an error when stop")
	}
}

func newBotMock() (*Bot, IConsumer, IChannel) {
	consumer := newMockConsumer()
	channel := newMockChannel()

	return NewBot(consumer, map[uuid.UUID]IChannel{
		channel.GetChannel().ChannelID: channel,
	}), consumer, channel
}
