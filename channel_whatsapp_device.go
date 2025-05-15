package lowbot

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type WhatsappDeviceChannel struct {
	*Channel
	running bool
	conn    *whatsmeow.Client
	device  *store.Device
	JID     *types.JID
}

var whatsSQLContainer *sqlstore.Container

func InitWhatsappDeviceChannel() {
	var err error
	address := os.Getenv("WHATSAPP_DEVICE_STORE_ADDRESS")

	if address == "" {
		address = "whatsapp_credentials.db?_foreign_keys=on"
	}

	go func() {
		whatsSQLContainer, err = sqlstore.New("sqlite3", address, nil)
	
		if err != nil {
			panic(err)
		}
	}()
}

func NewWhatsappDeviceChannel(JID *types.JID, callback func(whatsmeow.QRChannelItem, *types.JID) error) (IChannel, error) {
	var device *store.Device
	var conn *whatsmeow.Client
	var err error

	if JID == nil {
		device = whatsSQLContainer.NewDevice()
		conn = whatsmeow.NewClient(device, nil)

		qrChan, err := conn.GetQRChannel(context.Background())

		if err != nil {
			return nil, err
		}

		err = conn.Connect()

		if err != nil {
			return nil, err
		}

		for evt := range qrChan {
			err := callback(evt, device.ID)

			if err != nil {
				return nil, err
			}
		}

		conn.Disconnect()
	} else {
		device, err = whatsSQLContainer.GetDevice(*JID)

		if err != nil {
			return nil, err
		}

		conn = whatsmeow.NewClient(device, nil)
	}

	return &WhatsappDeviceChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_WHATSAPP_DEVICE_NAME,
			Broadcast: NewBroadcast[*Interaction](),
		},
		device:  device,
		conn:    conn,
		running: false,
	}, nil
}

func (channel *WhatsappDeviceChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *WhatsappDeviceChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	channel.conn.AddEventHandler(func(evt any) {
		switch v := evt.(type) {
		case *events.Message:
			{
				re := regexp.MustCompile(`:\d+`)
				rootUser := re.ReplaceAllString(v.Info.Sender.User, "")

				from := NewWho(rootUser, rootUser)
				from.Custom["JID"] = types.NewJID(rootUser, v.Info.Sender.Server)
				interaction := NewInteractionMessageText(v.Message.GetConversation())
				interaction.SetFrom(from)
				channel.Broadcast.Send(interaction)
			}
		}
	})

	err := channel.conn.Connect()

	if err != nil {
		return err
	}

	channel.running = true

	return nil
}

func (channel *WhatsappDeviceChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	err := channel.device.Delete()

	if err != nil {
		return err
	}

	err = channel.Broadcast.Close()

	if err != nil {
		return err
	}

	channel.conn.Disconnect()
	channel.running = false

	return nil
}

func (channel *WhatsappDeviceChannel) SendAudio(interaction *Interaction) error {
	err := channel.SendText(interaction)

	if err != nil {
		return err
	}

	err = interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(context.Background(), interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaAudio)

	if err != nil {
		return err
	}

	message := &waE2E.AudioMessage{
		Mimetype: proto.String(interaction.Parameters.File.GetFile().Mime),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}

	JID := interaction.From.Custom["JID"].(types.JID)

	_, err = channel.conn.SendMessage(context.Background(), JID, &waE2E.Message{
		AudioMessage: message,
	})

	return err
}

func (channel *WhatsappDeviceChannel) SendButton(interaction *Interaction) error {

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

func (channel *WhatsappDeviceChannel) SendDocument(interaction *Interaction) error {
	err := interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(context.Background(), interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaDocument)

	if err != nil {
		return err
	}

	message := &waE2E.DocumentMessage{
		Title:    proto.String(interaction.Parameters.File.GetFile().Name),
		Caption:  proto.String(interaction.Parameters.Text),
		Mimetype: proto.String(interaction.Parameters.File.GetFile().Mime),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}

	JID := interaction.From.Custom["JID"].(types.JID)

	_, err = channel.conn.SendMessage(context.Background(), JID, &waE2E.Message{
		DocumentMessage: message,
	})

	return err
}

func (channel *WhatsappDeviceChannel) SendImage(interaction *Interaction) error {
	err := interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(context.Background(), interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaImage)

	if err != nil {
		return err
	}

	message := &waE2E.ImageMessage{
		Caption:  proto.String(interaction.Parameters.Text),
		Mimetype: proto.String(interaction.Parameters.File.GetFile().Mime),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}

	JID := interaction.From.Custom["JID"].(types.JID)

	_, err = channel.conn.SendMessage(context.Background(), JID, &waE2E.Message{
		ImageMessage: message,
	})

	return err
}

func (channel *WhatsappDeviceChannel) SendText(interaction *Interaction) error {
	JID := interaction.From.Custom["JID"].(types.JID)
	_, err := channel.conn.SendMessage(context.Background(), JID, &waE2E.Message{
		Conversation: &interaction.Parameters.Text,
	})
	return err
}

func (channel *WhatsappDeviceChannel) SendVideo(interaction *Interaction) error {
	err := interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(context.Background(), interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaVideo)

	if err != nil {
		return err
	}

	message := &waE2E.VideoMessage{
		Caption:  proto.String(interaction.Parameters.Text),
		Mimetype: proto.String(interaction.Parameters.File.GetFile().Mime),

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}

	JID := interaction.From.Custom["JID"].(types.JID)

	_, err = channel.conn.SendMessage(context.Background(), JID, &waE2E.Message{
		VideoMessage: message,
	})

	return err
}
