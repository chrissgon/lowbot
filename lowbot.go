package lowbot

import (
	"log"
)

var Debug = true

func StartBot(base Flow, channel Channel, persist Persist) error {
	interactions := make(chan Interaction)

	go channel.Next(interactions)

	if Debug {
		log.Println("Bot is now running. Press CTRL-C to exit.")
	}

	for interaction := range interactions {
		flow, err := persist.Get(interaction.SessionID)

		if err != nil || flow.IsEnd() {
			flow = startFlow(interaction.SessionID, base)
		}

		processStep(flow, channel, interaction)

		persist.Set(flow)
	}

	close(interactions)

	return nil
}

func startFlow(sessionID string, flow Flow) *Flow {
	flow.Start()
	flow.SessionID = sessionID
	return &flow
}

func processStep(flow *Flow, channel Channel, in Interaction) error {
	next, err := RunAction(flow.Next(in), channel)

	if Debug {
		log.Printf("SessionID:<%v> Action:<%s> ERR:%v\n", in.SessionID, flow.Current.Action, err)
	}

	if err != nil {
		RunActionError(flow, channel)
		RunAction(flow.End(), channel)
		return NewError("RunAction", err)
	}

	if next {
		return processStep(flow, channel, in)
	}

	return err
}
