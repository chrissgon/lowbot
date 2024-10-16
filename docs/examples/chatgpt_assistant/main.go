package main

import (
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
)

func main() {
	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make consumer
	consumer, _ := lowbot.NewChatGPTAssistantConsumer(os.Getenv("CHATGPT_TOKEN"), os.Getenv("CHATGPT_ASSISTANT"))

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
