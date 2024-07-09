package main

import (
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

var userRoom *lowbot.Room

var fakeChannel = NewFakeChannel()
var fakeGuest = lowbot.NewWho("123", "fake guest")

func main() {

	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))
	// channel, _ := lowbot.NewDiscordChannel(os.Getenv("DISCORD_TOKEN"))

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

				// userRoom = room
				// userRoom.AddGuest(guest)
			}

			return true, nil
		},
	})

	// make consumer
	consumer := lowbot.NewJourneyConsumer(flow, persist)

	// start consumer
	lowbot.StartConsumer(consumer, []lowbot.IChannel{channel})
}
