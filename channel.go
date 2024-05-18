package lowbot

type Channel interface {
	SendAudio(*Interaction) error
	SendButton(*Interaction) error
	SendDocument(*Interaction) error
	SendImage(*Interaction) error
	SendText(*Interaction) error
	SendVideo(*Interaction) error
	Next(chan *Interaction)
}

type Interaction struct {
	SessionID  string                `json:"sessionID"`
	Type       InteractionType       `json:"type"`
	Parameters InteractionParameters `json:"parameters"`
	Custom     map[string]any        `json:"custom"`
}

type InteractionParameters StepParameters

type InteractionType string

const (
	MESSAGE InteractionType = "message"
	EVENT   InteractionType = "event"
)

func NewInteractionMessageAudio(sessionID string, audio string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE,
		Parameters: InteractionParameters{
			Text:  text,
			Audio: audio,
		},
	}
}

func NewInteractionMessageButton(sessionID string, buttons []string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE,
		Parameters: InteractionParameters{
			Text:    text,
			Buttons: buttons,
		},
	}
}

func NewInteractionMessageDocument(sessionID string, document string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE,
		Parameters: InteractionParameters{
			Text:     text,
			Document: document,
		},
	}
}

func NewInteractionMessageImage(sessionID string, image string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE,
		Parameters: InteractionParameters{
			Text:  text,
			Image: image,
		},
	}
}

func NewInteractionMessageText(sessionID string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE,
		Parameters: InteractionParameters{
			Text: text,
		},
	}
}

func NewInteractionMessageVideo(sessionID string, video string, text string) *Interaction {
	return &Interaction{
		SessionID: sessionID,
		Type:      MESSAGE,
		Parameters: InteractionParameters{
			Text:  text,
			Video: video,
		},
	}
}
