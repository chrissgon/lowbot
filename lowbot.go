package lowbot

import (
	"fmt"
)

var Debug = true
var AutoLoad = true

func StartBot(base Flow, channel Channel, persist Persist) error {
	ins := make(chan Interaction)

	go channel.Next(ins)

	for in := range ins {
		flow, err := persist.Get(in.SessionID)

		if err != nil || flow.IsEnd() {
			flow = startFlow(base)
		}

		err = runAction(flow, channel, in)

		if Debug {
			fmt.Printf("lowbot: %v\n", err)
		}

		persist.Set(in.SessionID, flow)
	}

	return nil
}

func startFlow(flow Flow) *Flow {
	flow.Start()
	return &flow
}

func runAction(flow *Flow, channel Channel, i Interaction) error {
	next, err := RunAction(i.SessionID, channel, flow.Next(i))

	if err != nil {
		RunActionError(i.SessionID, channel, flow)
		RunAction(i.SessionID, channel, flow.End())
		return NewError("RunAction", err)
	}

	if next {
		return runAction(flow, channel, i)
	}

	return err
}
