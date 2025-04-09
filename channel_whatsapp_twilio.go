package lowbot

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsappTwilioChannel struct {
	*Channel
	running bool
	conn    *twilio.RestClient
	server  *http.Server
	ctx     context.Context
	cancel  context.CancelFunc
	sid     string
}

func NewWhatsappTwilioChannel(token, SID string) (IChannel, error) {
	conn := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: SID,
		Password: token,
	})

	return &WhatsappTwilioChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_WHATSAPP_TWILIO_NAME,
			Broadcast: NewBroadcast[*Interaction](),
		},
		running: false,
		conn:    conn,
		sid:     SID,
	}, nil
}

func (channel *WhatsappTwilioChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *WhatsappTwilioChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.POST(fmt.Sprintf("/twilio/%v", channel.sid), func(c *gin.Context) {
		message := c.PostForm("Body")
		from := c.PostForm("From")
		to := c.PostForm("To")

		interaction := NewInteractionMessageText(message)

		interaction.SetFrom(NewWho(from, from))
		interaction.SetTo(NewWho(to, to))

		channel.Broadcast.Send(interaction)
	})

	port := os.Getenv("WHATSAPP_TWILIO_PORT")

	if port == "" {
		port = "8080"
	}

	channel.server = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	channel.ctx, channel.cancel = context.WithCancel(context.Background())

	go channel.server.ListenAndServe()

	channel.running = true

	return nil
}

func (channel *WhatsappTwilioChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	err := channel.server.Shutdown(channel.ctx)

	if err != nil {
		return err
	}
	
	err = channel.Broadcast.Close()
	
	if err != nil {
		return err
	}
	
	channel.cancel()
	channel.running = false

	return nil
}

func (channel *WhatsappTwilioChannel) SendAudio(interaction *Interaction) error {
	return channel.SendDocument(interaction)
}

func (channel *WhatsappTwilioChannel) SendButton(interaction *Interaction) error {
	_, contentSIDExists := interaction.Step.Parameters.Custom["contentSID"]
	_, contentVariablesExists := interaction.Step.Parameters.Custom["contentVariable"]

	if contentSIDExists && contentVariablesExists {
		contentSID, _ := interaction.Step.Parameters.Custom["contentSID"].(string)
		contentVariables, _ := interaction.Step.Parameters.Custom["contentVariable"].(string)

		to := interaction.From.WhoID
		from := interaction.To.WhoID

		params := &openapi.CreateMessageParams{}
		params.SetTo(to)
		params.SetFrom(from)
		params.SetBody(interaction.Parameters.Text)
		params.SetContentSid(contentSID)
		params.SetContentVariables(contentVariables)

		_, err := channel.conn.Api.CreateMessage(params)

		return err
	}

	sb := strings.Builder{}

	sb.WriteString(interaction.Parameters.Text)
	sb.WriteString("\n")

	for i, button := range interaction.Parameters.Buttons {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%v. %v", i+1, button))
	}

	interaction.Parameters.Text = sb.String()

	return channel.SendText(interaction)
}

func (channel *WhatsappTwilioChannel) SendDocument(interaction *Interaction) error {
	to := interaction.From.WhoID
	from := interaction.To.WhoID

	if !IsURL(interaction.Parameters.File.GetFile().Path) {
		return ERR_FILE_NOT_PUBLIC
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(interaction.Parameters.Text)
	params.SetMediaUrl([]string{
		interaction.Parameters.File.GetFile().Path,
	})

	_, err := channel.conn.Api.CreateMessage(params)

	return err
}

func (channel *WhatsappTwilioChannel) SendImage(interaction *Interaction) error {
	return channel.SendDocument(interaction)
}

func (channel *WhatsappTwilioChannel) SendText(interaction *Interaction) error {
	to := interaction.From.WhoID
	from := interaction.To.WhoID

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(interaction.Parameters.Text)

	_, err := channel.conn.Api.CreateMessage(params)

	return err
}

func (channel *WhatsappTwilioChannel) SendVideo(interaction *Interaction) error {
	return channel.SendDocument(interaction)
}
