package lowbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"time"
)

type ActionsMap map[string]ActionFunc

type ActionFunc func(*Flow, IChannel, Interaction) (bool, error)

var actions = ActionsMap{
	"Button":  RunActionButton,
	"File":    RunActionFile,
	"Input":   RunActionInput,
	"Text":    RunActionText,
	"Webhook": RunActionWebhook,
}

func SetCustomActions(custom ActionsMap) {
	maps.Copy(actions, custom)
}

func RunNextAction(flow *Flow, channel IChannel, interaction Interaction) (bool, error) {
	if flow == nil {
		return false, ERR_NIL_FLOW
	}

	action, exists := actions[flow.CurrentStep.Action]

	if !exists {
		return false, ERR_UNKNOWN_ACTION
	}

	return action(flow, channel, interaction)
}

func RunActionButton(flow *Flow, channel IChannel, interaction Interaction) (bool, error) {
	flow.Wait(interaction)
	defer func() {
		flow.Continue(interaction)
	}()

	step := flow.CurrentStep

	text := ParseTemplate(step.Parameters.Texts)

	answerInteraction := NewInteractionMessageButton(step.Parameters.Buttons, text)
	answerInteraction.SetFrom(interaction.From)
	answerInteraction.SetTo(interaction.To)
	answerInteraction.SetReplier(NewWho(flow.Replier, flow.Replier))

	err := SendInteraction(channel, answerInteraction)

	return false, err
}

func RunActionFile(flow *Flow, channel IChannel, interaction Interaction) (bool, error) {
	flow.Wait(interaction)
	defer func() {
		flow.Continue(interaction)
	}()

	step := flow.CurrentStep

	text := ParseTemplate(step.Parameters.Texts)

	fmt.Println("RunActionFile", step.Parameters.Path, step.Parameters.URL)
	answerInteraction := NewInteractionMessageFile(text, step.Parameters.Path, step.Parameters.URL)
	answerInteraction.SetFrom(interaction.From)
	answerInteraction.SetTo(interaction.To)
	answerInteraction.SetReplier(NewWho(flow.Replier, flow.Replier))

	err := SendInteraction(channel, answerInteraction)

	return true, err
}

func RunActionInput(flow *Flow, channel IChannel, interaction Interaction) (bool, error) {
	_, err := RunActionText(flow, channel, interaction)

	return false, err
}

func RunActionText(flow *Flow, channel IChannel, interaction Interaction) (bool, error) {
	flow.Wait(interaction)
	defer func() {
		flow.Continue(interaction)
	}()

	step := flow.CurrentStep

	text := ParseTemplate(step.Parameters.Texts)

	answerInteraction := NewInteractionMessageText(text)
	answerInteraction.SetFrom(interaction.From)
	answerInteraction.SetTo(interaction.To)
	answerInteraction.SetReplier(NewWho(flow.Replier, flow.Replier))

	err := SendInteraction(channel, answerInteraction)

	return true, err
}

func RunActionWebhook(flow *Flow, channel IChannel, interaction Interaction) (bool, error) {
	flow.Wait(interaction)
	defer func() {
		flow.Continue(interaction)
	}()

	step := flow.CurrentStep

	url := step.Parameters.URL
	timeout := step.Parameters.Timeout
	headers := step.Parameters.Headers

	body, err := json.Marshal(interaction)

	if err != nil {
		return true, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))

	if err != nil {
		return true, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	res, err := client.Do(req)

	if err != nil {
		return true, err
	}

	defer res.Body.Close()

	var answerInteractions []Interaction

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&answerInteractions)

	if err != nil {
		return true, err
	}

	for _, answerInteraction := range answerInteractions {
		answerInteraction.SetFrom(interaction.From)
		answerInteraction.SetTo(interaction.To)
		answerInteraction.SetReplier(NewWho(flow.Replier, flow.Replier))
		err := SendInteraction(channel, answerInteraction)

		if err != nil {
			return true, err
		}
	}

	return true, nil
}
