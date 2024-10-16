package main

import (
	"fmt"
	"os"
	"time"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

func main() {

	lowbot.DEBUG = true
	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make a persist
	persist, _ := lowbot.NewMemoryFlowPersist()

	// make consumer
	consumer := lowbot.NewJourneyConsumer(flow, persist)

	bot := lowbot.NewBot(consumer, map[uuid.UUID]lowbot.IChannel{
		channel.GetChannel().ChannelID: channel,
	})

	err := bot.Start()

	fmt.Println(err)

	// start consumer
	// ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(4 * time.Second)
		
		err := bot.Stop()
		fmt.Println("stopped", err)
		
		time.Sleep(4 * time.Second)
		err = bot.Start()
		fmt.Println("started", err)
		

		// consumer.GetConsumer().Cancel()
		// cancel()

	}()

	sc := make(chan os.Signal, 1)
	<-sc

	// lowbot.StartConsumer(consumer, []lowbot.IChannel{channel})

	// ctx.
}
