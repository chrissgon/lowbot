package main

import (
	"fmt"
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

func main() {
	lowbot.DEBUG = true

	// set custom actions
	lowbot.SetCustomActions(lowbot.ActionsMap{
		"TextUsername": func(interaction *lowbot.Interaction) (*lowbot.Interaction, bool) {
			template := lowbot.ParseTemplate(interaction.Step.Parameters.Texts)
			templateWithUsername := fmt.Sprintf(template, interaction.Parameters.Text)
			in := lowbot.NewInteractionMessageText(templateWithUsername)
			in.SetFrom(interaction.From)
			in.SetTo(interaction.To)
			return in, true
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
