package lowbot

import (
	"sync"
	"testing"
)

func TestStartBot(t *testing.T) {
	EnableLocalPersist = false
	SetCustomActions(ActionsMap{
		"Custom": func(flow *Flow, channel Channel) (bool, error) {
			return false, nil
		},
	})

	base, _ := NewFlow("./mocks/flow.yaml")
	// discord, _ := NewDiscord()
	telegram, _ := NewTelegram()
	persist, _ := NewLocalPersist()

	// go StartBot(base, discord, persist)
	go StartBot(base, telegram, persist)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
