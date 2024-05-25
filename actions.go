package lowbot

type ActionsMap map[string]func(flow *Flow, channel IChannel) (bool, error)

var actions = ActionsMap{
	"Audio":    ActionAudio,
	"Button":   ActionButton,
	"Document": ActionDocument,
	"Image":    ActionImage,
	"Input":    ActionInput,
	"Text":     ActionText,
	"Video":    ActionVideo,
	"Wait":     ActionWait,
}

func SetCustomActions(custom ActionsMap) {
	for k, v := range custom {
		actions[k] = v
	}
}

func RunAction(flow *Flow, channel IChannel) (bool, error) {
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

	return action(flow, channel)
}

func ActionAudio(flow *Flow, channel IChannel) (bool, error) {
	step := flow.Current
	err := channel.SendAudio(NewInteractionMessageAudio(channel.ChannelID(), flow.SessionID, step.Parameters.Audio, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func ActionButton(flow *Flow, channel IChannel) (bool, error) {
	step := flow.Current
	err := channel.SendButton(NewInteractionMessageButton(channel.ChannelID(), flow.SessionID, step.Parameters.Buttons, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func ActionDocument(flow *Flow, channel IChannel) (bool, error) {
	step := flow.Current
	err := channel.SendDocument(NewInteractionMessageDocument(channel.ChannelID(), flow.SessionID, step.Parameters.Document, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func ActionImage(flow *Flow, channel IChannel) (bool, error) {
	step := flow.Current
	err := channel.SendImage(NewInteractionMessageImage(channel.ChannelID(), flow.SessionID, step.Parameters.Image, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func ActionInput(flow *Flow, channel IChannel) (bool, error) {
	step := flow.Current
	err := channel.SendText(NewInteractionMessageText(channel.ChannelID(), flow.SessionID, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func ActionText(flow *Flow, channel IChannel) (bool, error) {
	step := flow.Current
	err := channel.SendText(NewInteractionMessageText(channel.ChannelID(), flow.SessionID, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func ActionVideo(flow *Flow, channel IChannel) (bool, error) {
	step := flow.Current
	err := channel.SendVideo(NewInteractionMessageVideo(channel.ChannelID(), flow.SessionID, step.Parameters.Video, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func ActionWait(flow *Flow, channel IChannel) (bool, error) {
	return true, nil
}
