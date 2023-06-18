package lowbot

import (
	"fmt"
)

type ActionsMap map[string]func(flow *Flow, channel Channel) (bool, error)

var actions = ActionsMap{
	"Audio":    actionAudio,
	"Button":   actionButton,
	"Document": actionDocument,
	"Image":    actionImage,
	"Input":    actionInput,
	"Text":     actionText,
	"Video":    actionVideo,
	"Wait":     actionWait,
}

func SetCustomActions(custom ActionsMap) {
	for k, v := range custom {
		actions[k] = v
	}
}

func runAction(flow *Flow, channel Channel) (bool, error) {
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

func runActionError(flow *Flow, channel Channel) (bool, error) {
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

func actionAudio(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendAudio(NewInteractionMessageAudio(flow.SessionID, step.Parameters.Audio, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func actionButton(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendButton(NewInteractionMessageButton(flow.SessionID, step.Parameters.Buttons, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func actionDocument(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendDocument(NewInteractionMessageDocument(flow.SessionID, step.Parameters.Document, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func actionImage(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendImage(NewInteractionMessageImage(flow.SessionID, step.Parameters.Image, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func actionInput(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendText(NewInteractionMessageText(flow.SessionID, ParseTemplate(step.Parameters.Texts)))
	return false, err
}

func actionText(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendText(NewInteractionMessageText(flow.SessionID, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func actionVideo(flow *Flow, channel Channel) (bool, error) {
	step := flow.Current
	err := channel.SendVideo(NewInteractionMessageVideo(flow.SessionID, step.Parameters.Video, ParseTemplate(step.Parameters.Texts)))
	return true, err
}

func actionWait(flow *Flow, channel Channel) (bool, error) {
	return false, nil
}
