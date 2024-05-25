package main

import (
	"fmt"
	"os"

	"github.com/chrissgon/lowbot"
)

func main() {
	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make a persist
	persist, _ := lowbot.NewMemoryFlowPersist()

	// set custom actions
	lowbot.SetCustomActions(lowbot.ActionsMap{
		"TextUsername": func(flow *lowbot.Flow, channel lowbot.IChannel) (bool, error) {
			step := flow.Current
			template := lowbot.ParseTemplate(step.Parameters.Texts)
			templateWithUsername := fmt.Sprintf(template, step.GetLastResponseText())
			in := lowbot.NewInteractionMessageText(channel.ChannelID(), flow.SessionID, templateWithUsername)
			err := channel.SendText(in)
			return true, err
		},
	})

	// make consumer
	consumer := lowbot.NewJourneyConsumer(flow, persist)

	// start consumer
	lowbot.StartConsumer(consumer, channel)
}
