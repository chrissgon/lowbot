package lowbot

import (
	"os"
	"testing"
)

func TestStartJourney(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	persist, _ := NewMemoryFlowPersist()
	consumer := NewJourneyConsumer(flow, persist)

	discord, _ := NewDiscordChannel(os.Getenv("DISCORD_TOKEN"))
	telegram, _ := NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	go StartConsumer(consumer, discord)
	StartConsumer(consumer, telegram)
}
