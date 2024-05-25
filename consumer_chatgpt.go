package lowbot

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTConsumer struct {
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
		conn:  conn,
		model: model,
	}, nil
}

func (c *ChatGPTConsumer) Run(interaction *Interaction, channel IChannel) {
	resp, _ := c.conn.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: interaction.Parameters.Text,
				},
			},
		},
	)

	channel.SendText(NewInteractionMessageText(channel.ChannelID(), interaction.SessionID, resp.Choices[0].Message.Content))
}
