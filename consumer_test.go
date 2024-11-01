package lowbot

import "github.com/google/uuid"

type mockConsumer struct {
	*Consumer
	ranTimes int
}

func newMockConsumer() IConsumer {
	return &mockConsumer{
		Consumer: &Consumer{
			ConsumerID: uuid.New(),
			Name:       "mock consumer",
		},
		ranTimes: 0,
	}
}

func (m *mockConsumer) GetConsumer() *Consumer {
	return m.Consumer
}

func (m *mockConsumer) Run(*Interaction) ([]*Interaction, error) {
	m.ranTimes++
	return nil, nil
}
