package lowbot

import "testing"

func TestActions_ActionsMap(t *testing.T) {
	names := []string{"Button", "File", "Input", "Text", "Wait"}

	for _, name := range names {
		_, exists := actions[name]

		if !exists {
			t.Errorf("not found action: %s", name)
		}
	}
}

func TestActions_SetCustomActions(t *testing.T) {
	name := "Custom"
	SetCustomActions(ActionsMap{
		name: func(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) { return true, nil },
	})

	_, exists := actions[name]

	if !exists {
		t.Errorf("not found action: %s", name)
	}
}
