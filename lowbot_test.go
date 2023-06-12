package lowbot

import "testing"

func TestStartBot(t *testing.T) {
	SetCustomActions(ActionsMap{
		"Custom": func(sessionID string, channel Channel, step *Step) (bool, error) {
			return false, nil
		},
	})

	base, _ := NewFlow("./mocks/button.yaml")
	channel, _ := NewDiscord()
	persist, _ := NewLocalPersist()

	go StartBot(base, channel, persist)
}
