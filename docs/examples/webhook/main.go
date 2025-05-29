package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	lowbot.DEBUG = true

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

	s := gin.Default()

	// create a server to receive the interactions
	s.POST("/", func(c *gin.Context) {
		var interaction lowbot.Interaction

		json.NewDecoder(c.Request.Body).Decode(&interaction)

		fmt.Println("received interaction", interaction.Parameters.Text)

		answerInteraction := lowbot.NewInteractionMessageText(fmt.Sprintf("Nice to meet you %v", interaction.Parameters.Text))
		answerInteraction.SetFrom(interaction.From)
		answerInteraction.SetTo(interaction.To)

		c.JSON(http.StatusOK, []lowbot.Interaction{answerInteraction})
	})

	s.Run(":3333")
}
