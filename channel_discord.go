package lowbot

import (
	"bytes"
	"os"
	"strings"

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

	return &DiscordChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_DISCORD_NAME,
		},
		conn: conn,
	}, nil
}

func (channel *DiscordChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *DiscordChannel) Next(interaction chan *Interaction) {
	channel.conn.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		sender := NewWho(m.ChannelID, s.State.User.Username)

		interaction <- NewInteractionMessageText(channel, sender, m.Content)
	})
	channel.conn.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		channel.RespondInteraction(i.Interaction)

		sender := NewWho(i.ChannelID, s.State.User.Username)

		interaction <- NewInteractionMessageText(channel, sender, i.Interaction.MessageComponentData().CustomID)
	})

	channel.conn.Identify.Intents = discordgo.IntentsGuildMessages

	err := channel.conn.Open()

	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	<-sc

	channel.conn.Close()
}

func (channel *DiscordChannel) RespondInteraction(in *discordgo.Interaction) {
	channel.conn.InteractionRespond(in, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: ""}},
	)
}

func (channel *DiscordChannel) SendAudio(interaction *Interaction) error {
	sessionID := interaction.Sender.WhoID
	path := interaction.Parameters.File.GetFile().Path

	return channel.SendFile(sessionID, interaction.Parameters.Text, path)
}

func (channel *DiscordChannel) SendButton(interaction *Interaction) error {
	sessionID := interaction.Sender.WhoID

	message := &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: channel.getButtons(interaction)},
		},
	}

	_, err := channel.conn.ChannelMessageSendComplex(sessionID, message)

	return err
}

func (*DiscordChannel) getButtons(interaction *Interaction) (buttons []discordgo.MessageComponent) {
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

func (channel *DiscordChannel) SendDocument(interaction *Interaction) error {
	sessionID := interaction.Sender.WhoID
	path := interaction.Parameters.File.GetFile().Path

	return channel.SendFile(sessionID, interaction.Parameters.Text, path)
}

func (channel *DiscordChannel) SendImage(interaction *Interaction) error {
	sessionID := interaction.Sender.WhoID
	path := interaction.Parameters.File.GetFile().Path

	if !IsURL(path) {
		return channel.SendFile(sessionID, interaction.Parameters.Text, path)
	}

	message := &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: path,
			},
		},
	}

	_, err := channel.conn.ChannelMessageSendComplex(sessionID, message)

	return err
}

func (channel *DiscordChannel) SendText(interaction *Interaction) error {
	sessionID := interaction.Sender.WhoID

	_, err := channel.conn.ChannelMessageSend(sessionID, interaction.Parameters.Text)

	return err
}

func (channel *DiscordChannel) SendVideo(interaction *Interaction) error {
	sessionID := interaction.Sender.WhoID
	path := interaction.Parameters.File.GetFile().Path

	if !IsURL(path) {
		return channel.SendFile(sessionID, interaction.Parameters.Text, path)
	}

	message := &discordgo.MessageSend{
		Content: interaction.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Video: &discordgo.MessageEmbedVideo{
				URL: path,
			},
		},
	}

	_, err := channel.conn.ChannelMessageSendComplex(sessionID, message)

	return err
}

func (channel *DiscordChannel) SendFile(sessionID, text, path string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]

	_, err = channel.conn.ChannelMessageSendComplex(sessionID, &discordgo.MessageSend{
		Content: text,
		File: &discordgo.File{
			Name:   name,
			Reader: bytes.NewReader(file),
		},
	})

	return err
}
