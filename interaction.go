package lowbot

import (
	"strings"
	"time"
)

type Interaction struct {
	// FlowID     uuid.UUID
	To      *Who
	From    *Who
	Replier *Who

	Type       InteractionType
	Parameters InteractionParameters
	Timestamp  int64
	Custom     map[string]any
}

type InteractionType string

const (
	MESSAGE_BUTTON InteractionType = "MESSAGE_BUTTON"
	MESSAGE_FILE   InteractionType = "MESSAGE_FILE"
	MESSAGE_TEXT   InteractionType = "MESSAGE_TEXT"
	EVENT_TYPING   InteractionType = "EVENT_TYPING"
)

type InteractionParameters struct {
	Buttons []string
	File    IFile
	Text    string

	Custom map[string]any
}

func NewInteractionMessageButton(buttons []string, text string) Interaction {
	return Interaction{
		Type:      MESSAGE_BUTTON,
		Timestamp: time.Now().Unix(),
		Parameters: InteractionParameters{
			Text:    text,
			Buttons: buttons,
		},
		Custom: map[string]any{},
	}
}

func NewInteractionMessageFile(text, path, url string) Interaction {
	return Interaction{
		Type:      MESSAGE_FILE,
		Timestamp: time.Now().Unix(),
		Parameters: InteractionParameters{
			Text: text,
			File: NewFile(path, url),
		},
		Custom: map[string]any{},
	}
}

func NewInteractionMessageText(text string) Interaction {
	return Interaction{
		Type:      MESSAGE_TEXT,
		Timestamp: time.Now().Unix(),
		Parameters: InteractionParameters{
			Text: text,
		},
		Custom: map[string]any{},
	}
}

func (interaction *Interaction) SetTo(to *Who) *Interaction {
	interaction.To = to
	return interaction
}

func (interaction *Interaction) SetFrom(from *Who) *Interaction {
	interaction.From = from
	return interaction
}

func (interaction *Interaction) SetReplier(replier *Who) *Interaction {
	interaction.Replier = replier
	return interaction
}

// func (interaction *Interaction) SetStep(step Step) *Interaction {
// 	interaction.Step = step
// 	return interaction
// }

func (interaction *Interaction) IsEmptyText() bool {
	return interaction.Parameters.Text == "" || strings.TrimSpace(interaction.Parameters.Text) == ""
}
