package lowbot

import (
	"testing"
)

func TestActionNames(t *testing.T) {
	names := []string{"Audio", "Button", "Document", "Image", "Input", "Text", "Video"}

	for _, name := range names {
		_, exists := actions[name]

		if !exists {
			t.Errorf("not found action: %s", name)
		}
	}
}

func TestSetCustomActions(t *testing.T) {
	name := "Custom"
	SetCustomActions(ActionsMap{
		name: func(sessionID string, channel Channel, step *Step) (bool, error) { return true, nil },
	})

	_, exists := actions[name]

	if !exists {
		t.Errorf("not found action: %s", name)
	}
}
