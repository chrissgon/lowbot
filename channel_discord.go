package lowbot

import (
	"bytes"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

type DiscordChannel struct {
	channelID uuid.UUID
	conn      *discordgo.Session
}

func NewDiscordChannel(token string) (IChannel, error) {
	if token == "" {
		return nil, ERR_UNKNOWN_DISCORD_TOKEN
	}

	conn, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	return &DiscordChannel{
		channelID: uuid.New(),
		conn:      conn,
	}, nil
}

func (ch *DiscordChannel) ChannelID() uuid.UUID {
	return ch.channelID
}

func (ch *DiscordChannel) Next(in chan *Interaction) {
	ch.conn.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		in <- NewInteractionMessageText(ch.channelID, m.ChannelID, m.Content)
	})
	ch.conn.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ch.RespondInteraction(i.Interaction)
		in <- NewInteractionMessageText(ch.channelID, i.ChannelID, i.Interaction.MessageComponentData().CustomID)
	})

	ch.conn.Identify.Intents = discordgo.IntentsGuildMessages

	err := ch.conn.Open()

	if err != nil {
		return
	}

	sc := make(chan os.Signal, 1)
	<-sc

	ch.conn.Close()
}

func (ch *DiscordChannel) RespondInteraction(in *discordgo.Interaction) {
	ch.conn.InteractionRespond(in, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: ""}},
	)
}

func (ch *DiscordChannel) SendAudio(interaction *Interaction) error {
	return ch.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Audio)
}

func (ch *DiscordChannel) SendButton(interaction *Interaction) error {
	_, err := ch.conn.ChannelMessageSendComplex(interaction.SessionID, &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: ch.getButtons(interaction)},
		},
	})
	return err
}

func (ch *DiscordChannel) getButtons(interaction *Interaction) (buttons []discordgo.MessageComponent) {
	for _, button := range interaction.Parameters.Buttons {
		buttons = append(buttons, discordgo.Button{
			Label:    button,
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			CustomID: button,
		})
	}
	return
}

func (ch *DiscordChannel) SendDocument(interaction *Interaction) error {
	return ch.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Document)
}

func (ch *DiscordChannel) SendImage(interaction *Interaction) error {
	if !IsURL(interaction.Parameters.Image) {
		return ch.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Image)
	}

	_, err := ch.conn.ChannelMessageSendComplex(interaction.SessionID, &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: interaction.Parameters.Image,
			},
		},
	})

	return err
}

func (ch *DiscordChannel) SendText(interaction *Interaction) error {
	_, err := ch.conn.ChannelMessageSend(interaction.SessionID, interaction.Parameters.Text)
	return err
}

func (ch *DiscordChannel) SendVideo(interaction *Interaction) error {
	if !IsURL(interaction.Parameters.Video) {
		return ch.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Video)
	}

	_, err := ch.conn.ChannelMessageSendComplex(interaction.SessionID, &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Video: &discordgo.MessageEmbedVideo{
				URL: interaction.Parameters.Video,
			},
		},
	})

	return err
}

func (ch *DiscordChannel) SendFile(sessionID, text, path string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]

	_, err = ch.conn.ChannelMessageSendComplex(sessionID, &discordgo.MessageSend{
		Content: text,
		File: &discordgo.File{
			Name:   name,
			Reader: bytes.NewReader(file),
		},
	})

	return err
}
