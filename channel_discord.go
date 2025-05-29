package lowbot

import (
	"bytes"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type DiscordChannel struct {
	*Channel
	conn *discordgo.Session
}

func NewDiscordChannel(token string) (IChannel, error) {
	if token == "" {
		return nil, ERR_UNKNOWN_DISCORD_TOKEN
	}

	conn, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	channel := &DiscordChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_DISCORD_NAME,
			Broadcast: NewBroadcast[Interaction](),
			Running:   false,
		},
		conn: conn,
	}

	return channel, nil
}

func (channel *DiscordChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *DiscordChannel) Start() error {
	if channel.Running {
		return ERR_CHANNEL_RUNNING
	}

	channel.conn.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !channel.Running {
			return
		}
		if m.Author.ID == s.State.User.ID {
			return
		}

		from := NewWho(m.ChannelID, s.State.User.Username)

		answerInteraction := NewInteractionMessageText(m.Content)
		answerInteraction.SetFrom(from)

		channel.Broadcast.Send(answerInteraction)

	})
	channel.conn.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		channel.RespondInteraction(i.Interaction)

		from := NewWho(i.ChannelID, s.State.User.Username)

		answerInteraction := NewInteractionMessageText(i.Interaction.MessageComponentData().CustomID)
		answerInteraction.SetFrom(from)

		channel.Broadcast.Send(answerInteraction)
	})

	channel.conn.Identify.Intents = discordgo.IntentsGuildMessages

	err := channel.conn.Open()

	if err != nil {
		channel.conn.Close()
		return err
	}

	channel.Running = true

	return nil
}

func (channel *DiscordChannel) Stop() error {
	if !channel.Running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	err := channel.Broadcast.Close()

	if err != nil {
		return err
	}

	err = channel.conn.Close()

	if err != nil {
		return err
	}

	channel.Running = false

	return nil
}

func (channel *DiscordChannel) RespondInteraction(in *discordgo.Interaction) {
	channel.conn.InteractionRespond(in, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: ""}},
	)
}

func (channel *DiscordChannel) SendAudio(interaction Interaction) error {
	sessionID := interaction.From.WhoID

	return channel.SendFile(sessionID, interaction)
}

func (channel *DiscordChannel) SendButton(interaction Interaction) error {
	sessionID := interaction.From.WhoID

	message := &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: channel.getButtons(interaction)},
		},
	}

	_, err := channel.conn.ChannelMessageSendComplex(sessionID, message)

	return err
}

func (*DiscordChannel) getButtons(interaction Interaction) (buttons []discordgo.MessageComponent) {
	for i, button := range interaction.Parameters.Buttons {
		buttons = append(buttons, discordgo.Button{
			Label:    button,
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			CustomID: strconv.Itoa(i + 1),
		})
	}
	return
}

func (channel *DiscordChannel) SendDocument(interaction Interaction) error {
	sessionID := interaction.From.WhoID
	return channel.SendFile(sessionID, interaction)
}

func (channel *DiscordChannel) SendImage(interaction Interaction) error {
	sessionID := interaction.From.WhoID
	return channel.SendFile(sessionID, interaction)
}

func (channel *DiscordChannel) SendText(interaction Interaction) error {
	if interaction.IsEmptyText() {
		return nil
	}

	sessionID := interaction.From.WhoID

	_, err := channel.conn.ChannelMessageSend(sessionID, interaction.Parameters.Text)

	return err
}

func (channel *DiscordChannel) SendVideo(interaction Interaction) error {
	sessionID := interaction.From.WhoID
	return channel.SendFile(sessionID, interaction)
}

func (channel *DiscordChannel) SendFile(sessionID string, interaction Interaction) error {
	err := interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	_, err = channel.conn.ChannelMessageSendComplex(sessionID, &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		File: &discordgo.File{
			Name:   interaction.Parameters.File.GetFile().Name,
			Reader: bytes.NewReader(interaction.Parameters.File.GetFile().Bytes),
		},
	})

	return err
}
