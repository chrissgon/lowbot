package main

import (
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// make a channel. In this exemple is Telegram
	channel, _ := lowbot.NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	// make consumer
	consumer, _ := lowbot.NewChatGPTConsumer(os.Getenv("CHATGPT_TOKEN"), openai.GPT3Dot5Turbo)

	// start bot
	lowbot.StartConsumer(consumer, channel)
}