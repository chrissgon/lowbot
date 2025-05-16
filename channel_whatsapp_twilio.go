package lowbot

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsappTwilioChannel struct {
	*Channel

	conn    *twilio.RestClient
	sid     string
	number  string
}

var whatsappTwilioCallbacks map[string]func(c *gin.Context) error = map[string]func(c *gin.Context) error{}
var whatsappTwilioCallbacksMutex = sync.Mutex{}

func InitWhatsappTwilioChannel(webhook *gin.Engine, path string) {
	webhook.POST(fmt.Sprintf("%s/:ID", path), func(c *gin.Context) {
		whatsappTwilioCallbacksMutex.Lock()
		defer whatsappTwilioCallbacksMutex.Unlock()

		ID := c.Param("ID")
		callback, exists := whatsappTwilioCallbacks[ID]

		if !exists {
			c.Status(http.StatusNotFound)
			return
		}

		err := callback(c)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Status(http.StatusOK)
	})
}

func NewWhatsappTwilioChannel(number, token, SID string) (IChannel, error) {
	conn := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: SID,
		Password: token,
	})

	return &WhatsappTwilioChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_WHATSAPP_TWILIO_NAME,
			Broadcast: NewBroadcast[*Interaction](),
			Running: false,
		},
		conn:    conn,
		sid:     SID,
		number:  number,
	}, nil
}

func (channel *WhatsappTwilioChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *WhatsappTwilioChannel) Start() error {
	if channel.Running {
		return ERR_CHANNEL_RUNNING
	}

	whatsappTwilioCallbacksMutex.Lock()
	defer whatsappTwilioCallbacksMutex.Unlock()

	whatsappTwilioCallbacks[channel.number] = func(c *gin.Context) error {
		message := c.PostForm("Body")
		from := c.PostForm("From")
		to := c.PostForm("To")

		interaction := NewInteractionMessageText(message)

		interaction.SetFrom(NewWho(from, from))
		interaction.SetTo(NewWho(to, to))

		channel.Broadcast.Send(interaction)

		return nil
	}

	channel.Running = true

	return nil
}

func (channel *WhatsappTwilioChannel) Stop() error {
	if !channel.Running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	err := channel.Broadcast.Close()

	if err != nil {
		return err
	}

	channel.Running = false

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
