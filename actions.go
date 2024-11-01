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

func GetAction(flow *Flow) (ActionFunc, error) {
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
	text := ParseTemplate(interaction.Step.Parameters.Texts)

	newInteraction := NewInteractionMessageButton(interaction.Step.Parameters.Buttons, text)
	newInteraction.SetFrom(interaction.From)
	newInteraction.SetTo(interaction.To)
	newInteraction.SetReplier(NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME))

	return newInteraction, true
}

func RunActionFile(interaction *Interaction) (*Interaction, bool) {
	text := ParseTemplate(interaction.Step.Parameters.Texts)

	newInteraction := NewInteractionMessageFile(interaction.Step.Parameters.Path, text)
	newInteraction.SetFrom(interaction.From)
	newInteraction.SetTo(interaction.To)
	newInteraction.SetReplier(NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME))

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
	text := ParseTemplate(interaction.Step.Parameters.Texts)

	newInteraction := NewInteractionMessageText(text)
	newInteraction.SetFrom(interaction.From)
	newInteraction.SetTo(interaction.To)
	newInteraction.SetReplier(NewWho(CONSUMER_JOURNEY_NAME, CONSUMER_JOURNEY_NAME))

	return newInteraction, false
}
