package main

import (
	"fmt"

	"github.com/chrissgon/lowbot"
)

func main() {
	// disable auto load persist
	lowbot.AutoLoad = false

	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegram()

	// make a persist
	persist, _ := lowbot.NewLocalPersist()

	// set custom actions
	lowbot.SetCustomActions(lowbot.ActionsMap{
		"TextUsername": func(sessionID string, channel lowbot.Channel, step *lowbot.Step) (bool, error) {
			template := lowbot.ParseTemplate(step.Parameters.Texts)
			templateWithUsername := fmt.Sprintf(template, step.GetLastResponseText())
			in := lowbot.NewInteractionMessageText(sessionID, templateWithUsername)
			err := channel.SendText(in)
			return lowbot.GetActionReturn(step.Action, true, err)
		},
	})

	// start bot
	lowbot.StartBot(flow, channel, persist)
}
