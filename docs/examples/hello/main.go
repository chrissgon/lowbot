package main

import (
	"fmt"
	"os"

	"github.com/chrissgon/lowbot"
)

func main() {
	// disable local persistence
	lowbot.EnableLocalPersist = false

	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegram(os.Getenv("TELEGRAM_TOKEN"))

	// make a persist
	persist, _ := lowbot.NewLocalPersist()

	// set custom actions
	lowbot.SetCustomActions(lowbot.ActionsMap{
		"TextUsername": func(flow *lowbot.Flow, channel lowbot.Channel) (bool, error) {
			step := flow.Current
			template := lowbot.ParseTemplate(step.Parameters.Texts)
			templateWithUsername := fmt.Sprintf(template, step.GetLastResponseText())
			in := lowbot.NewInteractionMessageText(flow.SessionID, templateWithUsername)
			err := channel.SendText(in)
			return true, err
		},
	})

	// start bot
	lowbot.StartBot(flow, channel, persist)
}
