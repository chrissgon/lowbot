package lowbot

import (
	"fmt"
	"time"
)

var Debug = true

func StartBot(base Flow, channel Channel, persist Persist) error {
	ins := make(chan Interaction)

	go channel.Next(ins)

	if Debug {
		fmt.Println("Bot is now running. Press CTRL-C to exit.")
	}

	for in := range ins {
		flow, err := persist.Get(in.SessionID)

		if err != nil || flow.IsEnd() {
			flow = startFlow(in.SessionID, base)
		}

		processStep(flow, channel, in)

		persist.Set(flow)
	}

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
		fmt.Printf("%v: <%v> Action%s %v\n", time.Now().UTC(), in.SessionID, flow.Current.Action, err)
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
