package lowbot

import "fmt"

type JourneyConsumer struct {
	flow    *Flow
	persist FlowPersist
}

func NewJourneyConsumer(flow *Flow, persist FlowPersist) IConsumer {
	return &JourneyConsumer{flow, persist}
}

func (journey *JourneyConsumer) Run(interaction *Interaction, channel IChannel) {
	flow, err := journey.persist.Get(interaction.SessionID)

	flowNotExistsOrWasFinished := err != nil || flow.NoHasNext()

	if flowNotExistsOrWasFinished {
		copyFlow := *journey.flow
		copyFlow.Start()
		copyFlow.SessionID = interaction.SessionID
		flow = &copyFlow
	}

	journey.processStep(flow, channel, interaction)

	journey.persist.Set(flow)
}

func (journey *JourneyConsumer) processStep(flow *Flow, channel IChannel, interaction *Interaction) error {
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

	return journey.processStep(flow, channel, interaction)
}
