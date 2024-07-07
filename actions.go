package lowbot

type ActionsMap map[string]func(*Flow, *Interaction, IChannel) (bool, error)

var actions = ActionsMap{
	"Button": RunActionButton,
	"File":   RunActionFile,
	"Input":  RunActionInput,
	"Text":   RunActionText,
	"Wait":   RunActionWait,
}

func SetCustomActions(custom ActionsMap) {
	for k, v := range custom {
		actions[k] = v
	}
}

func RunAction(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) {
	if flow == nil {
		return false, ERR_NIL_FLOW
	}

	step := flow.Current

	if step == nil {
		return false, ERR_NIL_STEP
	}

	action, exists := actions[step.Action]

	if !exists {
		return false, ERR_UNKNOWN_ACTION
	}

	return action(flow, interaction, channel)
}

func RunActionButton(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) {
	step := flow.Current

	text := ParseTemplate(step.Parameters.Texts)

	newInteraction := NewInteractionMessageButton(channel, interaction.Destination, interaction.Sender, step.Parameters.Buttons, text)

	return true, channel.SendButton(newInteraction)
}

func RunActionFile(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) {
	step := flow.Current

	text := ParseTemplate(step.Parameters.Texts)

	newInteraction := NewInteractionMessageFile(channel, interaction.Destination, interaction.Sender, step.Parameters.Path, text)

	if newInteraction.Parameters.File.IsAudio() {
		return false, channel.SendAudio(newInteraction)
	}
	if newInteraction.Parameters.File.IsImage() {
		return false, channel.SendImage(newInteraction)
	}
	if newInteraction.Parameters.File.IsVideo() {
		return false, channel.SendVideo(newInteraction)
	}

	return false, channel.SendDocument(newInteraction)
}

func RunActionInput(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) {
	_, err := RunActionText(flow, interaction, channel)

	return true, err
}

func RunActionText(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) {
	step := flow.Current

	text := ParseTemplate(step.Parameters.Texts)

	newInteraction := NewInteractionMessageText(channel, interaction.Destination, interaction.Sender, text)

	return false, channel.SendText(newInteraction)
}

func RunActionWait(flow *Flow, interaction *Interaction, channel IChannel) (bool, error) {
	return true, nil
}
