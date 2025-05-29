package main

import (
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

func main() {
	// Cria um fluxo
	flow, _ := lowbot.NewFlow("./flow.yaml")
	
	// Cria um canal (Usaremos o Telegram)
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// Cria um bot
	bot := lowbot.NewBot(flow, map[uuid.UUID]lowbot.IChannel{
		channel.GetChannel().ChannelID: channel,
	})

	// Inicia o bot
	bot.Start()

	// Mantem o processo rodando
	sc := make(chan os.Signal, 1)
	<-sc
}
