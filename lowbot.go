package lowbot

import (
	"fmt"
	"log"
)

var Debug = false

func StartBot(base *Flow, channel Channel, persist Persist) error {
	interactions := make(chan *Interaction)
	defer close(interactions)

	go channel.Next(interactions)

	if Debug {
		log.Println("Bot is now running. Press CTRL-C to exit.")
	}

	for interaction := range interactions {
		flow, err := persist.Get(interaction.SessionID)

		flowNotExistsOrWasFinished := err != nil || flow.NoHasNext()

		if flowNotExistsOrWasFinished {
			flow = startFlow(interaction.SessionID, *base)
		}

		processStep(flow, channel, interaction)

		persist.Set(flow)
	}

	return nil
}

func startFlow(sessionID string, flow Flow) *Flow {
	flow.Start()
	flow.SessionID = sessionID
	return &flow
}

func processStep(flow *Flow, channel Channel, interaction *Interaction) error {
	err := flow.Next(interaction)

	printLog(fmt.Sprintf("SessionID:<%v> Action:<%s> ERR: %v\n", interaction.SessionID, flow.Current.Action, err))

	if err != nil {
		return err
	}

	wait, err := RunAction(flow, channel)

	if err != nil {
		return err
	}

	if wait {
		return err
	}

	return processStep(flow, channel, interaction)
}

func printLog(msg string) {
	if Debug {
		log.Print(msg)
	}
}
