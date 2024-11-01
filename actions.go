package lowbot

type ActionsMap map[string]ActionFunc

type ActionFunc func(*Flow, *Interaction) (*Interaction, bool)

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

func RunActionButton(flow *Flow, interaction *Interaction) (*Interaction, bool) {
	step := flow.CurrentStep

	text := ParseTemplate(step.Parameters.Texts)

	newInteraction := NewInteractionMessageButton(interaction.Destination, NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME), step.Parameters.Buttons, text)

	return newInteraction, true
}

func RunActionFile(flow *Flow, interaction *Interaction) (*Interaction, bool) {
	step := flow.CurrentStep

	text := ParseTemplate(step.Parameters.Texts)

	newInteraction := NewInteractionMessageFile(interaction.Destination, NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME), step.Parameters.Path, text)

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

func RunActionInput(flow *Flow, interaction *Interaction) (*Interaction, bool) {
	answerInteraction, _ := RunActionText(flow, interaction)

	if answerInteraction.Parameters.Text == "" {
		return nil, true
	}

	return answerInteraction, true
}

func RunActionText(flow *Flow, interaction *Interaction) (*Interaction, bool) {
	step := flow.CurrentStep

	text := ParseTemplate(step.Parameters.Texts)

	newInteraction := NewInteractionMessageText(interaction.Destination, NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME), text)

	return newInteraction, false
}
