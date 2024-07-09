package lowbot

type Interaction struct {
	Channel     *Channel
	Destination *Who
	Sender      *Who
	Replier     *Who
	Type        InteractionType
	Parameters  InteractionParameters
	Custom      map[string]any
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
	Custom  map[string]any
}

func NewInteractionMessageButton(channel IChannel, destination *Who, sender *Who, buttons []string, text string) *Interaction {
	return &Interaction{
		Channel:     channel.GetChannel(),
		Destination: destination,
		Sender:      sender,
		Type:        MESSAGE_BUTTON,
		Parameters: InteractionParameters{
			Text:    text,
			Buttons: buttons,
		},
		Custom: map[string]any{},
	}
}

func NewInteractionMessageFile(channel IChannel, destination *Who, sender *Who, path string, text string) *Interaction {
	return &Interaction{
		Channel:     channel.GetChannel(),
		Destination: destination,
		Sender:      sender,
		Type:        MESSAGE_FILE,
		Parameters: InteractionParameters{
			Text: text,
			File: NewFile(path),
		},
		Custom: map[string]any{},
	}
}

func NewInteractionMessageText(channel IChannel, destination *Who, sender *Who, text string) *Interaction {
	return &Interaction{
		Channel:     channel.GetChannel(),
		Destination: destination,
		Sender:      sender,
		Type:        MESSAGE_TEXT,
		Parameters: InteractionParameters{
			Text: text,
		},
		Custom: map[string]any{},
	}
}

func (interaction *Interaction) SetReplier(replier *Who) *Interaction {
	interaction.Replier = replier
	return interaction
}

func (interaction *Interaction) SetDestination(destination *Who) *Interaction {
	interaction.Destination = destination
	return interaction
}
