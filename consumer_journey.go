package lowbot

import (
	"fmt"

	"github.com/google/uuid"
)

type JourneyConsumer struct {
	*Consumer
	Flow    *Flow
	Persist FlowPersist
}

func NewJourneyConsumer(flow *Flow, persist FlowPersist) IConsumer {
	return &JourneyConsumer{
		Consumer: &Consumer{
			ConsumerID: uuid.New(),
			Name:       CONSUMER_JOURNEY_NAME,
		},
		Flow:    flow,
		Persist: persist,
	}
}

func (consumer *JourneyConsumer) GetConsumer() *Consumer {
	return consumer.Consumer
}

func (consumer *JourneyConsumer) Run(interaction *Interaction, channel IChannel) error {
	flow, err := consumer.Persist.Get(interaction.Sender.WhoID)

	flowNotExistsOrWasFinished := err != nil || flow.NoHasNext()

	if flowNotExistsOrWasFinished {
		copyFlow := *consumer.Flow
		copyFlow.Start()
		flow = &copyFlow
	}

	err = consumer.processStep(flow, channel, interaction)

	consumer.Persist.Set(interaction.Sender.WhoID, flow)

	printLog(fmt.Sprintf("WhoID:<%v> Action:<%s> ERR: %v\n", interaction.Sender.WhoID, flow.Current.Action, err))

	return err
}

func (consumer *JourneyConsumer) processStep(flow *Flow, channel IChannel, interaction *Interaction) error {
	err := flow.Next(interaction)

	if err != nil {
		return err
	}

	replier := NewWho(flow.FlowID.String(), flow.Name)
	interaction.SetReplier(replier)

	wait, err := RunAction(flow, interaction, channel)

	if err != nil {
		return err
	}

	if wait {
		return nil
	}

	return consumer.processStep(flow, channel, interaction)
}
