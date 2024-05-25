package lowbot

import "testing"

func TestStartJourney(t *testing.T) {
	// os.Setenv("TELEGRAM_TOKEN", )
	// flow, _ := NewFlow("./mocks/flow.yaml")
	// persist, _ := NewMemoryFlowPersist()

	// joruney := NewJourneyConsumer(flow, persist)

	// lowbot.NewDiscordChannel("MTExNzU1Mjk0MTQxOTc0MTE4NA.G1_soj.kGkxobXrtE_Ko7DS7ODQOKnFPqCcNrS7IonNd0")
	// telegram, _ := NewTelegramChannel("7187351845:AAGGNkwB6ehdajFAJCrveKtyMe4mvMcO4Vg")
	// telegramB, _ := NewDiscordChannel("MTExNzU1Mjk0MTQxOTc0MTE4NA.G1_soj.kGkxobXrtE_Ko7DS7ODQOKnFPqCcNrS7IonNd0")

	// StartRoom([]IChannel{telegram, telegramB})
}

// func TestStartJourney(t *testing.T) {
// 	// os.Setenv("TELEGRAM_TOKEN", )
// 	flow, _ := NewFlow("./mocks/flow.yaml")
// 	persist, _ := NewMemoryFlowPersist()

// 	joruney := NewJourneyConsumer(flow, persist)

// 	channel, _ := NewDiscordChannel("MTExNzU1Mjk0MTQxOTc0MTE4NA.G1_soj.kGkxobXrtE_Ko7DS7ODQOKnFPqCcNrS7IonNd0")
// 	// telegramA, _ := NewTelegramChannel("7187351845:AAGGNkwB6ehdajFAJCrveKtyMe4mvMcO4Vg")
// 	// telegramB, _ := NewTelegramChannel("6824630284:AAFbdBJK5bSSjkjqfFGLiroxwUawhBl7AVk")

// 	StartConsumer(joruney, channel)
// }

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
