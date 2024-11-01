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
	text := ParseTemplate(interaction.StepParameters.Texts)

	newInteraction := NewInteractionMessageButton(interaction.StepParameters.Buttons, text)
	newInteraction.SetFrom(interaction.From)
	newInteraction.SetTo(interaction.To)
	newInteraction.SetReplier(NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME))

	return newInteraction, true
}

func RunActionFile(interaction *Interaction) (*Interaction, bool) {
	text := ParseTemplate(interaction.StepParameters.Texts)

	newInteraction := NewInteractionMessageFile(interaction.StepParameters.Path, text)
	newInteraction.SetFrom(interaction.From)
	newInteraction.SetTo(interaction.To)
	newInteraction.SetReplier(NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME))

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
	text := ParseTemplate(interaction.StepParameters.Texts)

	newInteraction := NewInteractionMessageText(text)
	newInteraction.SetFrom(interaction.From)
	newInteraction.SetTo(interaction.To)
	newInteraction.SetReplier(NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME))

	return newInteraction, false
}
