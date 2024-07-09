package main

import (
	"fmt"
	"time"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

type FakeChannel struct {
	*lowbot.Channel
}

func NewFakeChannel() lowbot.IChannel {
	channel := &FakeChannel{
		Channel: &lowbot.Channel{
			ChannelID: uuid.New(),
			Name:      "fake channel",
			Broadcast: lowbot.NewBroadcast[*lowbot.Interaction](),
		},
	}
	go channel.Next()
	return channel
}

// Close implements lowbot.IChannel.
func (f *FakeChannel) Close() error {
	panic("unimplemented")
}

// GetChannel implements lowbot.IChannel.
func (f *FakeChannel) GetChannel() *lowbot.Channel {
	return f.Channel
}

// Next implements lowbot.IChannel.
func (f *FakeChannel) Next() {
	for {
		time.Sleep(3 * time.Second)
		f.Channel.Broadcast.Send(lowbot.NewInteractionMessageText(
			f,
			fakeGuest,
			fakeGuest,
			"Fake automatic message",
		))
	}
}

// SendAudio implements lowbot.IChannel.
func (f *FakeChannel) SendAudio(*lowbot.Interaction) error {
	panic("unimplemented")
}

// SendButton implements lowbot.IChannel.
func (f *FakeChannel) SendButton(*lowbot.Interaction) error {
	panic("unimplemented")
}

// SendDocument implements lowbot.IChannel.
func (f *FakeChannel) SendDocument(*lowbot.Interaction) error {
	panic("unimplemented")
}

// SendImage implements lowbot.IChannel.
func (f *FakeChannel) SendImage(*lowbot.Interaction) error {
	panic("unimplemented")
}

// SendText implements lowbot.IChannel.
func (f *FakeChannel) SendText(interaction *lowbot.Interaction) error {
	fmt.Printf("User %s said: %s \n", interaction.Sender.Name, interaction.Parameters.Text)

	return nil
}

// SendVideo implements lowbot.IChannel.
func (f *FakeChannel) SendVideo(*lowbot.Interaction) error {
	panic("unimplemented")
}
