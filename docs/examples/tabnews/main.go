package main

import (
	"os"

	"github.com/chrissgon/lowbot"
)

func main() {
	// Cria um fluxo
	flow, _ := lowbot.NewFlow("./flow.yaml")
	// Cria um canal (Usaremos o Telegram)
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))
	// Cria uma persistência
	persist, _ := lowbot.NewMemoryFlowPersist()

	// make consumer
	consumer := lowbot.NewJourneyConsumer(flow, persist)

	// start consumer
	lowbot.StartConsumer(consumer, []lowbot.IChannel{channel})
}
