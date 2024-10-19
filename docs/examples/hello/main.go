package main

import (
	"fmt"
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
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
		"TextUsername": func(flow *lowbot.Flow, interaction *lowbot.Interaction) (*lowbot.Interaction, bool) {
			step := flow.CurrentStep
			template := lowbot.ParseTemplate(step.Parameters.Texts)
			templateWithUsername := fmt.Sprintf(template, flow.GetLastResponseText())
			in := lowbot.NewInteractionMessageText(interaction.Destination, interaction.Sender, templateWithUsername)
			// err := channel.SendText(in)
			return in, true
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
