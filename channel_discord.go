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
	conn   *discordgo.Session
	running bool
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
			Broadcast: NewBroadcast[*Interaction](),
		},
		conn:   conn,
		running: false,
	}

	return channel, nil
}

func (channel *DiscordChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *DiscordChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	channel.conn.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !channel.running {
			return
		}
		if m.Author.ID == s.State.User.ID {
			return
		}

		destination := NewWho(m.ChannelID, s.State.User.Username)

		channel.Broadcast.Send(NewInteractionMessageText(channel, destination, destination, m.Content))
	})
	channel.conn.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		channel.RespondInteraction(i.Interaction)

		destination := NewWho(i.ChannelID, s.State.User.Username)

		channel.Broadcast.Send(NewInteractionMessageText(channel, destination, destination, i.Interaction.MessageComponentData().CustomID))
	})

	channel.conn.Identify.Intents = discordgo.IntentsGuildMessages

	err := channel.conn.Open()

	if err != nil {
		channel.conn.Close()
		return err
	}

	channel.running = true

	return nil
}

func (channel *DiscordChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	channel.running = false
	channel.Broadcast.Close()
	return channel.conn.Close()
}

func (channel *DiscordChannel) RespondInteraction(in *discordgo.Interaction) {
	channel.conn.InteractionRespond(in, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: ""}},
	)
}

func (channel *DiscordChannel) SendAudio(interaction *Interaction) error {
	sessionID := interaction.Destination.WhoID
	path := interaction.Parameters.File.GetFile().Path

	return channel.SendFile(sessionID, interaction.Parameters.Text, path)
}

func (channel *DiscordChannel) SendButton(interaction *Interaction) error {
	sessionID := interaction.Destination.WhoID

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
	sessionID := interaction.Destination.WhoID
	path := interaction.Parameters.File.GetFile().Path

	return channel.SendFile(sessionID, interaction.Parameters.Text, path)
}

func (channel *DiscordChannel) SendImage(interaction *Interaction) error {
	sessionID := interaction.Destination.WhoID
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
	sessionID := interaction.Destination.WhoID

	_, err := channel.conn.ChannelMessageSend(sessionID, interaction.Parameters.Text)

	return err
}

func (channel *DiscordChannel) SendVideo(interaction *Interaction) error {
	sessionID := interaction.Destination.WhoID
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
