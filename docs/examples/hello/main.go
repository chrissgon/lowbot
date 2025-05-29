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
		"TextUsername": func(flow *lowbot.Flow, channel lowbot.IChannel, interaction lowbot.Interaction) (bool, error) {
			template := lowbot.ParseTemplate(flow.CurrentStep.Parameters.Texts)
			templateWithUsername := fmt.Sprintf(template, interaction.Parameters.Text)
			in := lowbot.NewInteractionMessageText(templateWithUsername)
			in.SetFrom(interaction.From)
			in.SetTo(interaction.To)

			err := lowbot.SendInteraction(channel, in)
			return false, err
		},
	})

	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make bot
	bot := lowbot.NewBot(flow, map[uuid.UUID]lowbot.IChannel{
		channel.GetChannel().ChannelID: channel,
	})

	// start bot
	bot.Start()

	// keep the process running
	sc := make(chan os.Signal, 1)
	<-sc
}
