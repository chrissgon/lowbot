package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

type Bot struct {
	BotID    uuid.UUID
	Flow     *Flow
	Channels map[uuid.UUID]IChannel
	Running  bool
}

func NewBot(flow *Flow, channels map[uuid.UUID]IChannel) *Bot {
	return &Bot{
		BotID:    uuid.New(),
		Flow:     flow,
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

		go bot.StartListenChannel(channel)
	}

	bot.Running = true

	return nil
}

func (bot *Bot) StartChannel(channel IChannel) error {
	return channel.Start()
}

func (bot *Bot) StartListenChannel(channel IChannel) {
	listener := channel.GetChannel().Broadcast.Listen()

	for interaction := range listener {
		flow, err := FlowPersist.Get(interaction.From.WhoID)

		flowNotExistsOrWasFinished := err != nil || flow.Ended()

		if flowNotExistsOrWasFinished {
			copyFlow := *bot.Flow
			copyFlow.Start()
			flow = &copyFlow
		}

		if flow.Waiting {
			continue
		}
		next := true

		for next {
			err := flow.Next(interaction)

			if err != nil {
				PrintLog(fmt.Sprintf("Channel:<%v> WhoID:<%v> ERR: %v\n", channel.GetChannel().Name, interaction.From.WhoID, err))
				break
			}

			next, err = RunNextAction(flow, channel, interaction)

			if err == nil {
				FlowPersist.Set(interaction.From.WhoID, flow)
				continue
			}

			PrintLog(fmt.Sprintf("Channel:<%v> WhoID:<%v> ERR: %v\n", channel.GetChannel().Name, interaction.From.WhoID, err))

			flow.NextError()

			next, err = RunNextAction(flow, channel, interaction)

			if err != nil {
				PrintLog(fmt.Sprintf("Channel:<%v> WhoID:<%v> ERR: %v\n", channel.GetChannel().Name, interaction.From.WhoID, err))
			}

			FlowPersist.Set(interaction.From.WhoID, flow)
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
