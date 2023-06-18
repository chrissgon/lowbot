package lowbot

import (
	"fmt"
)

type ActionsMap map[string]func(flow *Flow, channel Channel) (bool, error)

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

func RunAction(flow *Flow, channel Channel) (bool, error) {
	if flow == nil {
		return false, nil
	}

	step := flow.Current

	if step == nil {
		return false, nil
	}

	action, exists := actions[step.Action]

	if !exists {
		return false, NewError("RunAction", fmt.Errorf("not found action: %s", step.Action))
	}

	return action(flow, channel)
}

func RunActionError(flow *Flow, channel Channel) (bool, error) {
	if flow == nil {
		return false, NewError("RunActionError", fmt.Errorf("nil flow"))
	}

	step := flow.Steps["error"]

	if step == nil {
		return false, NewError("RunActionError", fmt.Errorf("not found step: error"))
	}

	action, exists := actions[step.Action]

	if !exists {
		return false, NewError("RunActionError", fmt.Errorf("not found action: %s", step.Action))
	}

	return action(flow, channel)
}

func ActionAudio(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendAudio(NewInteractionMessageAudio(flow.SessionID, step.Parameters.Audio, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func ActionButton(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendButton(NewInteractionMessageButton(flow.SessionID, step.Parameters.Buttons, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func ActionDocument(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendDocument(NewInteractionMessageDocument(flow.SessionID, step.Parameters.Document, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func ActionImage(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendImage(NewInteractionMessageImage(flow.SessionID, step.Parameters.Image, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func ActionInput(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendText(NewInteractionMessageText(flow.SessionID, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func ActionText(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendText(NewInteractionMessageText(flow.SessionID, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func ActionVideo(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendVideo(NewInteractionMessageVideo(flow.SessionID, step.Parameters.Video, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func ActionWait(flow *Flow, channel Channel) (bool, error) {
	return false, nil
}
