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

// GetConsumer implements IConsumer.
func (m *mockConsumer) GetConsumer() *Consumer {
	return m.Consumer
}

// Run implements IConsumer.
func (m *mockConsumer) Run(*Interaction, IChannel) error {
	m.ranTimes++
	return nil
}
