package lowbot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type WhatsappMeowChannel struct {
	*Channel
	conn      *whatsmeow.Client
	device    *store.Device
	JID       *types.JID
	handlerID uint32

	ctx    context.Context
	cancel context.CancelFunc
}

var whatsMeowSQLContainer *sqlstore.Container
var WhatsMeowReplyDuration time.Duration = 2 * time.Second

func InitWhatsappMeowChannel(ctx context.Context, db *sql.DB, dialect string, log waLog.Logger) error {
	whatsMeowSQLContainer = sqlstore.NewWithDB(db, dialect, log)

	return whatsMeowSQLContainer.Upgrade(ctx)
}

func NewWhatsappMeowChannel(JID *types.JID, callback func(whatsmeow.QRChannelItem, *types.JID) error) (IChannel, error) {
	var device *store.Device
	var conn *whatsmeow.Client
	var err error

	if JID == nil {
		device = whatsMeowSQLContainer.NewDevice()
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
	} else {
		device, err = whatsMeowSQLContainer.GetDevice(context.Background(), *JID)

		if err != nil {
			return nil, err
		}

		if device == nil {
			return nil, errors.New("unknown device")
		}

		conn = whatsmeow.NewClient(device, nil)

		conn.Connect()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WhatsappMeowChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_WHATSAPP_DEVICE_NAME,
			Broadcast: NewBroadcast[*Interaction](),
			Running:   false,
		},
		ctx:       ctx,
		cancel:    cancel,
		device:    device,
		conn:      conn,
		handlerID: 0,
	}, nil
}

func (channel *WhatsappMeowChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *WhatsappMeowChannel) Start() error {
	if channel.Running {
		return ERR_CHANNEL_RUNNING
	}

	channel.handlerID = channel.conn.AddEventHandler(func(evt any) {
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

	channel.Running = true

	return nil
}

func (channel *WhatsappMeowChannel) Stop() error {
	if !channel.Running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	err := channel.Broadcast.Close()

	if err != nil {
		return err
	}

	channel.conn.RemoveEventHandler(channel.handlerID)
	channel.Running = false

	return nil
}

func (channel *WhatsappMeowChannel) SendAudio(interaction *Interaction) error {
	err := channel.SendText(interaction)

	if err != nil {
		return err
	}

	err = interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(channel.ctx, interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaAudio)

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

	// sleep to look like a human reply
	time.Sleep(WhatsMeowReplyDuration)

	_, err = channel.conn.SendMessage(channel.ctx, JID, &waE2E.Message{
		AudioMessage: message,
	})

	return err
}

func (channel *WhatsappMeowChannel) SendButton(interaction *Interaction) error {
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

func (channel *WhatsappMeowChannel) SendDocument(interaction *Interaction) error {
	err := interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(channel.ctx, interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaDocument)

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

	// sleep to look like a human reply
	time.Sleep(WhatsMeowReplyDuration)

	_, err = channel.conn.SendMessage(channel.ctx, JID, &waE2E.Message{
		DocumentMessage: message,
	})

	return err
}

func (channel *WhatsappMeowChannel) SendImage(interaction *Interaction) error {
	err := interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(channel.ctx, interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaImage)

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

	// sleep to look like a human reply
	time.Sleep(WhatsMeowReplyDuration)

	_, err = channel.conn.SendMessage(channel.ctx, JID, &waE2E.Message{
		ImageMessage: message,
	})

	return err
}

func (channel *WhatsappMeowChannel) SendText(interaction *Interaction) error {
	JID := interaction.From.Custom["JID"].(types.JID)

	// sleep to look like a human reply
	time.Sleep(WhatsMeowReplyDuration)

	_, err := channel.conn.SendMessage(channel.ctx, JID, &waE2E.Message{
		Conversation: &interaction.Parameters.Text,
	})

	return err
}

func (channel *WhatsappMeowChannel) SendVideo(interaction *Interaction) error {
	err := interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	resp, err := channel.conn.Upload(channel.ctx, interaction.Parameters.File.GetFile().Bytes, whatsmeow.MediaVideo)

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

	// sleep to look like a human reply
	time.Sleep(WhatsMeowReplyDuration)

	_, err = channel.conn.SendMessage(channel.ctx, JID, &waE2E.Message{
		VideoMessage: message,
	})

	return err
}
