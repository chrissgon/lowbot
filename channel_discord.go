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

func (ds *DiscordChannel) ChannelID() uuid.UUID {
	return ds.channelID
}

func (ds *DiscordChannel) Next(in chan *Interaction) {
	ds.conn.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		in <- NewInteractionMessageText(m.ChannelID, m.Content)
	})
	ds.conn.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ds.RespondInteraction(i.Interaction)
		in <- NewInteractionMessageText(i.ChannelID, i.Interaction.MessageComponentData().CustomID)
	})

	ds.conn.Identify.Intents = discordgo.IntentsGuildMessages

	err := ds.conn.Open()

	if err != nil {
		return
	}

	sc := make(chan os.Signal, 1)
	<-sc

	ds.conn.Close()
}

func (ds *DiscordChannel) RespondInteraction(in *discordgo.Interaction) {
	ds.conn.InteractionRespond(in, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: ""}},
	)
}

func (ds *DiscordChannel) SendAudio(interaction *Interaction) error {
	return ds.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Audio)
}

func (ds *DiscordChannel) SendButton(interaction *Interaction) error {
	_, err := ds.conn.ChannelMessageSendComplex(interaction.SessionID, &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: ds.getButtons(interaction)},
		},
	})
	return err
}

func (ds *DiscordChannel) getButtons(interaction *Interaction) (buttons []discordgo.MessageComponent) {
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

func (ds *DiscordChannel) SendDocument(interaction *Interaction) error {
	return ds.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Document)
}

func (ds *DiscordChannel) SendImage(interaction *Interaction) error {
	if !IsURL(interaction.Parameters.Image) {
		return ds.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Image)
	}

	_, err := ds.conn.ChannelMessageSendComplex(interaction.SessionID, &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: interaction.Parameters.Image,
			},
		},
	})
	return err
}

func (ds *DiscordChannel) SendText(interaction *Interaction) error {
	_, err := ds.conn.ChannelMessageSend(interaction.SessionID, interaction.Parameters.Text)
	return err
}

func (ds *DiscordChannel) SendVideo(interaction *Interaction) error {
	if !IsURL(interaction.Parameters.Video) {
		return ds.SendFile(interaction.SessionID, interaction.Parameters.Text, interaction.Parameters.Video)
	}

	_, err := ds.conn.ChannelMessageSendComplex(interaction.SessionID, &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Video: &discordgo.MessageEmbedVideo{
				URL: interaction.Parameters.Video,
			},
		},
	})

	return err
}

func (ds *DiscordChannel) SendFile(sessionID, text, path string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]

	_, err = ds.conn.ChannelMessageSendComplex(sessionID, &discordgo.MessageSend{
		Content: text,
		File: &discordgo.File{
			Name:   name,
			Reader: bytes.NewReader(file),
		},
	})

	return err
}
