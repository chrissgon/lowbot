package lowbot

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
)

type ChatGPTConsumer struct {
	*Consumer
	model string
	conn  *openai.Client
}

type ChatGPTAssistantConsumer struct {
	*Consumer
	conn      *openai.Client
	assistant openai.Assistant
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

func NewChatGPTAssistantConsumer(token string, assistantID string) (IConsumer, error) {
	if token == "" {
		return nil, ERR_UNKNOWN_CHATGPT_TOKEN
	}
	if assistantID == "" {
		return nil, ERR_UNDEFINED_CHATGPT_ASSISTANT
	}

	conn := openai.NewClient(token)

	assistant, err := conn.RetrieveAssistant(context.Background(), assistantID)

	if err != nil {
		return nil, err
	}

	return &ChatGPTAssistantConsumer{
		Consumer: &Consumer{
			ConsumerID: uuid.New(),
			Name:       CONSUMER_CHATGPT_NAME,
		},
		conn:      conn,
		assistant: assistant,
	}, nil
}

func (consumer *ChatGPTConsumer) Run(interaction *Interaction, channel IChannel) error {
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
		printLog(fmt.Sprintf("%v: WhoID:<%v> ERR: %v\n", consumer.Name, interaction.Sender.WhoID, err))
		return err
	}

	replier := NewWho(consumer.ConsumerID, consumer.Name)
	newInteraction := NewInteractionMessageText(channel, interaction.Sender, resp.Choices[0].Message.Content)
	newInteraction.SetReplier(replier)

	return channel.SendText(newInteraction)
}

func (consumer *ChatGPTAssistantConsumer) Run(interaction *Interaction, channel IChannel) error {
	run, err := consumer.conn.CreateThreadAndRun(context.Background(), openai.CreateThreadAndRunRequest{
		RunRequest: openai.RunRequest{
			AssistantID:  consumer.assistant.ID,
			Model:        consumer.assistant.Model,
			Instructions: *consumer.assistant.Instructions,
		},
		Thread: openai.ThreadRequest{
			Messages: []openai.ThreadMessage{
				{
					Role:    openai.ThreadMessageRoleUser,
					Content: interaction.Parameters.Text,
				},
			},
		},
	})

	if err != nil {
		printLog(fmt.Sprintf("%v: WhoID:<%v> ERR: %v\n", consumer.Name, interaction.Sender.WhoID, err))
		return err
	}

	replier := NewWho(consumer.ConsumerID, consumer.Name)
	newInteraction := NewInteractionMessageText(channel, interaction.Sender, consumer.waitMessage(run))
	newInteraction.SetReplier(replier)

	return channel.SendText(newInteraction)
}

func (consumer *ChatGPTAssistantConsumer) waitMessage(run openai.Run) string {
	if run.Status == "completed" {
		limit := 1
		order := "desc"
		after := ""
		before := ""

		msgs, _ := consumer.conn.ListMessage(context.Background(), run.ThreadID, &limit, &order, &after, &before)

		return msgs.Messages[0].Content[0].Text.Value
	}

	time.Sleep(2 * time.Second)
	run, _ = consumer.conn.RetrieveRun(context.Background(), run.ThreadID, run.ID)
	return consumer.waitMessage(run)
}
