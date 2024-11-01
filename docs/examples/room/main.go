package main

import (
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

var fakeChannel = NewFakeChannel()
var fakeGuest = lowbot.NewWho("123", "fake guest")

func main() {

	// set custom action
	lowbot.SetCustomActions(lowbot.ActionsMap{
		"Fakechat": func(flow *lowbot.Flow, interaction *lowbot.Interaction) (*lowbot.Interaction, bool) {

			return nil, true
		},
	})

	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make a persist
	persist, _ := lowbot.NewMemoryFlowPersist()

	// make consumer
	consumer := lowbot.NewJourneyConsumer(flow, persist)

	// make bot
	bot := lowbot.NewBot(consumer, map[uuid.UUID]lowbot.IChannel{
		channel.GetChannel().ChannelID: channel,
	})

	// start bot
	bot.Start()

	// keep the process running
	sc := make(chan os.Signal, 1)
	<-sc
}
