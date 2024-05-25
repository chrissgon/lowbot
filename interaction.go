package lowbot

type Interaction struct {
	Channel *Channel
	Sender  *Who
	Replier *Who
	Type       InteractionType
	Parameters InteractionParameters
	Custom     map[string]any
}

type InteractionType string

const (
	MESSAGE_BUTTON InteractionType = "button"
	MESSAGE_FILE   InteractionType = "file"
	MESSAGE_TEXT   InteractionType = "text"
	EVENT_TYPING   InteractionType = "typing"
)

type InteractionParameters struct {
	Buttons []string
	File    IFile
	Text    string
	Custom  map[string]any
}

func NewInteractionMessageButton(channel IChannel, sender *Who, buttons []string, text string) *Interaction {
	return &Interaction{
		Channel: channel.GetChannel(),
		Sender:  sender,
		Type:    MESSAGE_BUTTON,
		Parameters: InteractionParameters{
			Text:    text,
			Buttons: buttons,
		},
	}
}

func NewInteractionMessageFile(channel IChannel, sender *Who, path string, text string) *Interaction {
	return &Interaction{
		Channel: channel.GetChannel(),
		Sender:  sender,
		Type:    MESSAGE_FILE,
		Parameters: InteractionParameters{
			Text: text,
			File: NewFile(path),
		},
	}
}

func NewInteractionMessageText(channel IChannel, sender *Who, text string) *Interaction {
	return &Interaction{
		Channel: channel.GetChannel(),
		Sender:  sender,
		Type:    MESSAGE_TEXT,
		Parameters: InteractionParameters{
			Text: text,
		},
	}
}

func (interaction *Interaction) SetReplier(replier *Who) *Interaction {
	interaction.Replier = replier
	return interaction
}
