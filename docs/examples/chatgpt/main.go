package main

import (
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

func main() {
	lowbot.DEBUG = true
	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make consumer
	consumer, _ := lowbot.NewChatGPTConsumer(os.Getenv("CHATGPT_TOKEN"), openai.GPT4Turbo)

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
