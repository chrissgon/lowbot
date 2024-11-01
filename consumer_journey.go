package lowbot

import (
	"fmt"
	"maps"

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

func (consumer *JourneyConsumer) Run(interaction *Interaction) ([]*Interaction, error) {
	flow, err := consumer.Persist.Get(interaction.Sender.WhoID)

	flowNotExistsOrWasFinished := err != nil || flow.Ended()

	if flowNotExistsOrWasFinished {
		copyFlow := *consumer.Flow
		copyFlow.Start()
		flow = &copyFlow
	}

	interactions, err := consumer.getInteractions(flow, interaction)

	consumer.Persist.Set(interaction.Sender.WhoID, flow)

	printLog(fmt.Sprintf("WhoID:<%v> Step:<%s> ERR: %v\n", interaction.Sender.WhoID, flow.CurrentStepName, err))

	return interactions, err
}

func (consumer *JourneyConsumer) getInteractions(flow *Flow, interaction *Interaction) ([]*Interaction, error) {
	next := true
	interactions := []*Interaction{}

	for next {
		err := flow.Next(interaction)

		if err != nil {
			return interactions, err
		}

		action, err := GetAction(flow, interaction)

		if err != nil {
			return interactions, err
		}

		answerInteraction, wait := action(flow, interaction)

		maps.Copy(flow.CurrentStep.Parameters.Custom, answerInteraction.Parameters.Custom)

		next = !wait

		if answerInteraction != nil {
			interactions = append(interactions, answerInteraction)
		}
	}

	return interactions, nil
}
