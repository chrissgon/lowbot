package lowbot

import (
	"fmt"
)

type ActionsMap map[string]func(sessionID string, channel Channel, step *Step) (bool, error)

var actions = ActionsMap{
	"Audio":    ActionAudio,
	"Button":   ActionButton,
	"Document": ActionDocument,
	"Image":    ActionImage,
	"Input":    ActionInput,
	"Text":     ActionText,
	"Video":    ActionVideo,
}

func SetCustomActions(custom ActionsMap) {
	for k, v := range custom {
		actions[k] = v
	}
}

func RunAction(sessionID string, channel Channel, step *Step) (bool, error) {
	if step == nil {
		return false, nil
	}

	action, exists := actions[step.Action]

	if !exists {
		return false, NewError("RunAction", fmt.Errorf("not found action: %s", step.Action))
	}

	return action(sessionID, channel, step)
}

func RunActionError(sessionID string, channel Channel, flow *Flow) (bool, error) {
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

	return action(sessionID, channel, step)
}

func GetActionReturn(name string, next bool, err error) (bool, error) {
	if err != nil {
		return next, NewError("Action"+name, err)
	}
	return next, nil
}

func ActionAudio(sessionID string, channel Channel, step *Step) (bool, error) {
	err := channel.SendAudio(NewInteractionMessageAudio(sessionID, step.Parameters.Audio, ParseTemplate(step.Parameters.Texts)))
	return GetActionReturn(step.Action, true, err)
}

func ActionButton(sessionID string, channel Channel, step *Step) (bool, error) {
	err := channel.SendButton(NewInteractionMessageButton(sessionID, step.Parameters.Buttons, ParseTemplate(step.Parameters.Texts)))
	return GetActionReturn(step.Action, false, err)
}

func ActionDocument(sessionID string, channel Channel, step *Step) (bool, error) {
	err := channel.SendDocument(NewInteractionMessageDocument(sessionID, step.Parameters.Document, ParseTemplate(step.Parameters.Texts)))
	return GetActionReturn(step.Action, true, err)
}

func ActionImage(sessionID string, channel Channel, step *Step) (bool, error) {
	err := channel.SendImage(NewInteractionMessageImage(sessionID, step.Parameters.Image, ParseTemplate(step.Parameters.Texts)))
	return GetActionReturn(step.Action, true, err)
}

func ActionInput(sessionID string, channel Channel, step *Step) (bool, error) {
	err := channel.SendText(NewInteractionMessageText(sessionID, ParseTemplate(step.Parameters.Texts)))
	return GetActionReturn(step.Action, false, err)
}

func ActionText(sessionID string, channel Channel, step *Step) (bool, error) {
	err := channel.SendText(NewInteractionMessageText(sessionID, ParseTemplate(step.Parameters.Texts)))
	return GetActionReturn(step.Action, true, err)
}

func ActionVideo(sessionID string, channel Channel, step *Step) (bool, error) {
	err := channel.SendVideo(NewInteractionMessageVideo(sessionID, step.Parameters.Video, ParseTemplate(step.Parameters.Texts)))
	return GetActionReturn(step.Action, true, err)
}
