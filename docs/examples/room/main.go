package main

import (
	"fmt"
	"os"

	"github.com/chrissgon/lowbot"
)

var fakeChannel = NewFakeChannel()
var fakeGuest = lowbot.NewWho("123", "fake guest")

func main() {
	
	
	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")
	
	// make a channel. In this exemple is Telegram
	telegramChannel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make a persist
	persist, _ := lowbot.NewMemoryFlowPersist()

	// make consumer
	consumer := lowbot.NewJourneyConsumer(flow, persist)

	journeyChannel, _ := lowbot.NewConsumerChannel(consumer)

	guests := map[string]*lowbot.Guest{
		lowbot.CHANNEL_TELEGRAM_NAME: lowbot.NewGuest(lowbot.NewWho(lowbot.CHANNEL_TELEGRAM_NAME, lowbot.CHANNEL_TELEGRAM_NAME), telegramChannel),
		lowbot.CONSUMER_JOURNEY_NAME: lowbot.NewGuest(lowbot.NewWho(lowbot.CONSUMER_JOURNEY_NAME, lowbot.CONSUMER_JOURNEY_NAME), journeyChannel),
	}

	room := lowbot.NewRoom(guests)

	err := room.Start()

	fmt.Println(err)



	// // set custom action
	// lowbot.SetCustomActions(lowbot.ActionsMap{
	// 	"Fakechat": func(flow *lowbot.Flow, interaction *lowbot.Interaction) (*lowbot.Interaction, bool) {
	// 		roomID, exists := flow.GetLastResponse().Custom["RoomID"].(uuid.UUID)

	// 		if exists {
	// 			roomManager := lowbot.GetRoomManager()
	// 			guest := lowbot.NewGuest(fakeGuest, fakeChannel)

	// 			roomManager.AddGuest(roomID, guest)
	// 		}

	// 		return nil, true
	// 	},
	// })

	// // make a flow
	// flow, _ := lowbot.NewFlow("./flow.yaml")

	// // make a channel. In this exemple is Telegram
	// channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// // make a persist
	// persist, _ := lowbot.NewMemoryFlowPersist()

	// // make consumer
	// consumer := lowbot.NewJourneyConsumer(flow, persist)

	// // make bot
	// bot := lowbot.NewBot(consumer, map[uuid.UUID]lowbot.IChannel{
	// 	channel.GetChannel().ChannelID: channel,
	// })

	// // start bot
	// bot.Start()

	// keep the process running
	sc := make(chan os.Signal, 1)
	<-sc
}
