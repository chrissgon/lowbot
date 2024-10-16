package main

import (
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

var fakeChannel = NewFakeChannel()
var fakeGuest = lowbot.NewWho("123", "fake guest")

func main() {
	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make a persist
	persist, _ := lowbot.NewMemoryFlowPersist()

	// set custom action
	lowbot.SetCustomActions(lowbot.ActionsMap{
		"Fakechat": func(flow *lowbot.Flow, interaction *lowbot.Interaction, channel lowbot.IChannel) (bool, error) {
			roomID, exists := flow.GetLastResponse().Custom["RoomID"].(uuid.UUID)

			if exists {
				roomManager := lowbot.GetRoomManager()
				guest := lowbot.NewGuest(fakeGuest, fakeChannel)

				roomManager.AddGuest(roomID, guest)
			}

			return true, nil
		},
	})

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
