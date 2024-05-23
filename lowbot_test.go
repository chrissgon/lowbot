package lowbot

import (
	"os"
	"testing"
)

func TestStartJourney(t *testing.T) {
	os.Setenv("TELEGRAM_TOKEN", "7187351845:AAGGNkwB6ehdajFAJCrveKtyMe4mvMcO4Vg")
	flow, _ := NewFlow("./mocks/flow.yaml")
	persist, _ := NewMemoryFlowPersist()

	joruney := NewJourneyConsumer(flow, persist)

	telegram, _ := NewTelegramChannel(os.Getenv("TELEGRAM_TOKEN"))

	StartConsumer(joruney, telegram)
}

// func TestStartBot(t *testing.T) {
// 	EnableLocalPersist = false
// 	SetCustomActions(ActionsMap{
// 		"Custom": func(flow *Flow, channel Channel) (bool, error) {
// 			return false, nil
// 		},
// 	})

// 	base, _ := NewFlow("./mocks/flow.yaml")
// 	// discord, _ := NewDiscord(os.Getenv("DISCORD_TOKEN"))
// 	telegram, err := NewTelegram(os.Getenv("TELEGRAM_TOKEN"))
// 	persist, _ := NewLocalPersist()

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// go StartBot(base, discord, persist)
// 	go StartBot(base, telegram, persist)

// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	wg.Wait()
// }
