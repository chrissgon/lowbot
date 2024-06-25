package lowbot

import (
	"context"
	"errors"
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
	threads   map[any]string
	ctx       context.Context
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
	ctx := context.Background()

	assistant, err := conn.RetrieveAssistant(ctx, assistantID)

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
		threads:   map[any]string{},
		ctx:       ctx,
	}, nil
}

func (consumer *ChatGPTConsumer) GetConsumer() *Consumer {
	return consumer.Consumer
}

func (consumer *ChatGPTAssistantConsumer) GetConsumer() *Consumer {
	return consumer.Consumer
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

	replier := NewWho(consumer.ConsumerID.String(), consumer.Name)
	newInteraction := NewInteractionMessageText(channel, interaction.Sender, resp.Choices[0].Message.Content)
	newInteraction.SetReplier(replier)

	return channel.SendText(newInteraction)
}

func (consumer *ChatGPTAssistantConsumer) Run(interaction *Interaction, channel IChannel) error {
	threadID, err := consumer.getThreadID(interaction)

	consumer.threads[interaction.Sender.WhoID] = threadID

	if err != nil {
		return err
	}

	run, err := consumer.conn.CreateRun(consumer.ctx, threadID, openai.RunRequest{
		AssistantID:  consumer.assistant.ID,
		Model:        consumer.assistant.Model,
		Instructions: *consumer.assistant.Instructions,
	})

	if err != nil {
		return err
	}

	answer, err := consumer.waitMessage(run)

	if err != nil {
		return err
	}

	replier := NewWho(consumer.ConsumerID.String(), consumer.Name)
	newInteraction := NewInteractionMessageText(channel, interaction.Sender, answer)
	newInteraction.SetReplier(replier)

	return channel.SendText(newInteraction)
}

func (consumer *ChatGPTAssistantConsumer) getThreadID(interaction *Interaction) (string, error) {
	threadID, exists := consumer.threads[interaction.Sender.WhoID]

	if exists {
		consumer.conn.CreateMessage(consumer.ctx, threadID, openai.MessageRequest{
			Role:    string(openai.ThreadMessageRoleUser),
			Content: interaction.Parameters.Text,
		})

		return threadID, nil
	}

	thread, err := consumer.conn.CreateThread(consumer.ctx, openai.ThreadRequest{
		Messages: []openai.ThreadMessage{
			{
				Role:    openai.ThreadMessageRoleUser,
				Content: interaction.Parameters.Text,
			},
		},
	})

	return thread.ID, err
}

func (consumer *ChatGPTAssistantConsumer) waitMessage(run openai.Run) (string, error) {
	if run.Status == openai.RunStatusCompleted {
		limit := 1
		order := "desc"
		after := ""
		before := ""

		msgs, err := consumer.conn.ListMessage(consumer.ctx, run.ThreadID, &limit, &order, &after, &before)

		if err != nil {
			return "", err
		}

		return msgs.Messages[0].Content[0].Text.Value, nil
	}
	if run.Status == "queued" || run.Status == "in_progress" {
		time.Sleep(1 * time.Second)

		run, err := consumer.conn.RetrieveRun(consumer.ctx, run.ThreadID, run.ID)

		if err != nil {
			return "", err
		}

		return consumer.waitMessage(run)
	}

	return "", errors.New(run.LastError.Message)
}
