package lowbot

import (
	"context"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type ChatGPTConsumer struct {
	*Consumer
	model string
	conn  *openai.Client
}

func NewChatGPTConsumer(token string, model string) (IConsumer, error) {
	if token == "" {
		return nil, ERR_UNKNOWN_CHATGPT_TOKEN
	}

	conn := openai.NewClient(token)

	if conn == nil {
		return nil, ERR_CONNECT_CHATGPT
	}

	return &ChatGPTConsumer{
		Consumer: &Consumer{
			ConsumerID: uuid.New(),
			Name:       CONSUMER_CHATGPT_NAME,
		},
		conn:  conn,
		model: model,
	}, nil
}

func (consumer *ChatGPTConsumer) GetConsumer() *Consumer {
	return consumer.Consumer
}

func (consumer *ChatGPTConsumer) Run(interaction *Interaction) ([]*Interaction, error) {
	resp, err := consumer.conn.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: consumer.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: interaction.Parameters.Text,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	answerInteraction := NewInteractionMessageText(resp.Choices[0].Message.Content)

	replier := NewWho(consumer.ConsumerID.String(), consumer.Name)
	answerInteraction.SetReplier(replier)
	answerInteraction.SetTo(interaction.To)
	answerInteraction.SetFrom(interaction.From)

	return []*Interaction{answerInteraction}, nil
}
