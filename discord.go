package lowbot

import (
	"bytes"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	conn *discordgo.Session
}

func (ds *Discord) Next(in chan *Interaction) {
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

func (ds *Discord) RespondInteraction(in *discordgo.Interaction) {
	ds.conn.InteractionRespond(in, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: ""}},
	)
}

func (ds *Discord) SendAudio(in *Interaction) error {
	return ds.SendFile(in.SessionID, in.Parameters.Text, in.Parameters.Audio)
}

func (ds *Discord) SendButton(in *Interaction) error {
	_, err := ds.conn.ChannelMessageSendComplex(in.SessionID, &discordgo.MessageSend{
		Content: in.Parameters.Text,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: ds.getButtons(in)},
		},
	})
	return err
}

func (ds *Discord) getButtons(in *Interaction) (buttons []discordgo.MessageComponent) {
	for _, button := range in.Parameters.Buttons {
		buttons = append(buttons, discordgo.Button{
			Label:    button,
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			CustomID: button,
		})
	}
	return
}

func (ds *Discord) SendDocument(in *Interaction) error {
	return ds.SendFile(in.SessionID, in.Parameters.Text, in.Parameters.Document)
}

func (ds *Discord) SendImage(in *Interaction) error {
	if !IsURL(in.Parameters.Image) {
		return ds.SendFile(in.SessionID, in.Parameters.Text, in.Parameters.Image)
	}

	_, err := ds.conn.ChannelMessageSendComplex(in.SessionID, &discordgo.MessageSend{
		Content: in.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: in.Parameters.Image,
			},
		},
	})
	return err
}

func (ds *Discord) SendText(in *Interaction) error {
	_, err := ds.conn.ChannelMessageSend(in.SessionID, in.Parameters.Text)
	return err
}

func (ds *Discord) SendVideo(in *Interaction) error {
	if !IsURL(in.Parameters.Video) {
		return ds.SendFile(in.SessionID, in.Parameters.Text, in.Parameters.Video)
	}

	_, err := ds.conn.ChannelMessageSendComplex(in.SessionID, &discordgo.MessageSend{
		Content: in.Parameters.Text,
		Embed: &discordgo.MessageEmbed{
			Video: &discordgo.MessageEmbedVideo{
				URL: in.Parameters.Video,
			},
		},
	})

	return err
}

func (ds *Discord) SendFile(sessionID, text, path string) error {
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

func NewDiscord() (Channel, error) {
	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		return nil, NewError("NewDiscord", ERR_UNKNOWN_DISCORD_TOKEN)
	}

	conn, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	return &Discord{conn: conn}, nil
}
