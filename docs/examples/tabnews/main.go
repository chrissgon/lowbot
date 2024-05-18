package main

import (
	"os"

	"github.com/chrissgon/lowbot"
)

func main() {
	// Desabilita persistência local
	lowbot.EnableLocalPersist = false

	// Cria um fluxo
	flow, _ := lowbot.NewFlow("./flow.yaml")
	// Cria um canal (Usaremos o Telegram)
	channel, _ := lowbot.NewTelegram(os.Getenv("TELEGRAM_TOKEN"))
	// Cria uma persistência
	persist, _ := lowbot.NewLocalPersist()

	// Inicia o bot
	lowbot.StartBot(flow, channel, persist)
}
