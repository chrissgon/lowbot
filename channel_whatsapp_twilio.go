package lowbot

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsappTwilioChannel struct {
	*Channel

	conn   *twilio.RestClient
	sid    string
	number string
}

type WhatsappTwilioNumbersManager struct {
	callback   func(c *gin.Context) error
	statusChan chan string
}

var whatsappTwilioNumbersManager map[string]WhatsappTwilioNumbersManager = map[string]WhatsappTwilioNumbersManager{}
var whatsappTwilioCallbacksMutex = sync.Mutex{}
var WhatsappTwilioReplyDuration time.Duration = 4 * time.Second
var WhatsappTwilioDefaultPath string = "/runner/whatsapp/twilio"

func InitWhatsappTwilioChannel(webhook *gin.Engine, path string) {
	WhatsappTwilioDefaultPath = path

	webhook.POST(fmt.Sprintf("%s/:ID", WhatsappTwilioDefaultPath), func(c *gin.Context) {
		whatsappTwilioCallbacksMutex.Lock()
		defer whatsappTwilioCallbacksMutex.Unlock()

		ID := c.Param("ID")
		manager, exists := whatsappTwilioNumbersManager[ID]

		if !exists {
			c.Status(http.StatusNotFound)
			return
		}

		err := manager.callback(c)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.Status(http.StatusOK)
	})

	go func() {
		webhook.POST(fmt.Sprintf("%s/:ID/status", WhatsappTwilioDefaultPath), func(c *gin.Context) {
			whatsappTwilioCallbacksMutex.Lock()
			defer whatsappTwilioCallbacksMutex.Unlock()

			ID := c.Param("ID")
			manager, exists := whatsappTwilioNumbersManager[ID]

			if !exists {
				c.Status(http.StatusNotFound)
				return
			}

			status := c.PostForm("MessageStatus")

			if status == "delivered" {
				manager.statusChan <- status
			}

			c.Status(http.StatusOK)
		})
	}()
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
			Broadcast: NewBroadcast[Interaction](),
			Running:   false,
		},
		conn:   conn,
		sid:    SID,
		number: number,
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

	whatsappTwilioNumbersManager[channel.number] = WhatsappTwilioNumbersManager{
		statusChan: make(chan string),

		callback: func(c *gin.Context) error {
			message := c.PostForm("Body")
			from := c.PostForm("From")
			to := c.PostForm("To")

			interaction := NewInteractionMessageText(message)

			interaction.SetFrom(NewWho(from, from))
			interaction.SetTo(NewWho(to, to))

			channel.Broadcast.Send(interaction)

			return nil
		},
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

func (channel *WhatsappTwilioChannel) SendAudio(interaction Interaction) error {
	return channel.SendDocument(interaction)
}

func (channel *WhatsappTwilioChannel) SendButton(interaction Interaction) error {
	step, err := GetCurrentStep(interaction.From.WhoID)

	if err != nil {
		return err
	}

	_, contentSIDExists := step.Parameters.Custom["contentSID"]
	_, contentVariablesExists := step.Parameters.Custom["contentVariable"]

	if contentSIDExists && contentVariablesExists {
		contentSID, _ := step.Parameters.Custom["contentSID"].(string)
		contentVariables, _ := step.Parameters.Custom["contentVariable"].(string)

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

func (channel *WhatsappTwilioChannel) SendDocument(interaction Interaction) error {
	err := channel.SendText(interaction)

	if err != nil {
		return err
	}

	to := interaction.From.WhoID
	from := interaction.To.WhoID

	if !IsURL(interaction.Parameters.File.GetFile().URL) {
		return ERR_FILE_NOT_PUBLIC
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetMediaUrl([]string{
		interaction.Parameters.File.GetFile().URL,
	})

	_, err = channel.conn.Api.CreateMessage(params)

	if err != nil {
		return err
	}

	return channel.getMessageStatus(interaction)
}

func (channel *WhatsappTwilioChannel) SendImage(interaction Interaction) error {
	return channel.SendDocument(interaction)
}

func (channel *WhatsappTwilioChannel) SendText(interaction Interaction) error {
	if interaction.IsEmptyText() {
		return nil
	}

	to := interaction.From.WhoID
	from := interaction.To.WhoID

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(interaction.Parameters.Text)

	_, err := channel.conn.Api.CreateMessage(params)

	if err != nil {
		return err
	}

	return channel.getMessageStatus(interaction)
}

func (channel *WhatsappTwilioChannel) SendVideo(interaction Interaction) error {
	return channel.SendDocument(interaction)
}

func (channel *WhatsappTwilioChannel) getMessageStatus(interaction Interaction) error {
	from := interaction.To.WhoID

	ID := strings.Split(from, ":")[1]
	manager, exists := whatsappTwilioNumbersManager[ID]

	if !exists {
		return fmt.Errorf("unknown number")
	}

	status := <-manager.statusChan

	if status != "delivered" {
		return fmt.Errorf("undelivered message")
	}

	return nil
}
