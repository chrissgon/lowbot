package lowbot

type ActionsMap map[string]ActionFunc

type ActionFunc func(*Interaction) (*Interaction, bool)

var actions = ActionsMap{
	"Button": RunActionButton,
	"File":   RunActionFile,
	"Input":  RunActionInput,
	"Text":   RunActionText,
}

func SetCustomActions(custom ActionsMap) {
	for k, v := range custom {
		actions[k] = v
	}
}

func GetAction(flow *Flow, interaction *Interaction) (ActionFunc, error) {
	if flow == nil {
		return nil, ERR_NIL_FLOW
	}

	step := flow.CurrentStep

	if step == nil {
		return nil, ERR_NIL_STEP
	}

	action, exists := actions[step.Action]

	if !exists {
		return nil, ERR_UNKNOWN_ACTION
	}

	return action, nil
}

func RunActionButton(interaction *Interaction) (*Interaction, bool) {
	newInteraction := NewInteractionMessageButton(interaction.Destination, NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME), interaction.Custom["buttons"].([]string), interaction.Custom["text"].(string))

	return newInteraction, true
}

func RunActionFile(interaction *Interaction) (*Interaction, bool) {
	newInteraction := NewInteractionMessageFile(interaction.Destination, NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME), interaction.Custom["path"].(string), interaction.Custom["text"].(string))

	if newInteraction.Parameters.File.IsAudio() {
		return newInteraction, false
	}
	if newInteraction.Parameters.File.IsImage() {
		return newInteraction, false
	}
	if newInteraction.Parameters.File.IsVideo() {
		return newInteraction, false
	}

	return newInteraction, false
}

func RunActionInput(interaction *Interaction) (*Interaction, bool) {
	answerInteraction, _ := RunActionText(interaction)

	if answerInteraction.Parameters.Text == "" {
		return nil, true
	}

	return answerInteraction, true
}

func RunActionText(interaction *Interaction) (*Interaction, bool) {
	newInteraction := NewInteractionMessageText(interaction.Destination, NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME), interaction.Custom["text"].(string))

	return newInteraction, false
}
