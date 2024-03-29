package main

import (
	"github.com/chrissgon/lowbot"
)

func main() {
	// disable local persistence
	lowbot.EnableLocalPersist = false

	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegram()

	// make a persist
	persist, _ := lowbot.NewLocalPersist()

	// start bot
	lowbot.StartBot(flow, channel, persist)
}
