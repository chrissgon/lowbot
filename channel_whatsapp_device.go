package lowbot

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type WhatsappChannel struct {
	*Channel
	running bool
	conn    *whatsmeow.Client
}

func NewWhatsappDeviceChannel() (IChannel, error) {
	return &WhatsappChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_TELEGRAM_NAME,
			Broadcast: NewBroadcast[*Interaction](),
		},
		running: false,
	}, nil
}

func (channel *WhatsappChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *WhatsappChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	go func() {
		container, err := sqlstore.New("sqlite3", "../file:whatsapp_credentials.db?_foreign_keys=on", nil)
		if err != nil {
			panic(err)
		}
		// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
		deviceStore, err := container.GetFirstDevice()
		if err != nil {
			panic(err)
		}
		// client := whatsmeow.NewClient(deviceStore, nil)

		// clientLog := waLog.Stdout("Client", "DEBUG", true)
		channel.conn = whatsmeow.NewClient(deviceStore, nil)
		channel.conn.AddEventHandler(func(evt interface{}) {
			switch v := evt.(type) {
			case *events.Message:
				{
					fmt.Println("Received a message!", v.Message.GetConversation())
					// fmt.Println(v.Info.Sender.id)

					destination := NewWho(v.Info.Sender.User, v.Info.Sender.User)
					destination.Custom["JID"] = v.Info.Sender
					interaction := NewInteractionMessageText(destination, destination, v.Message.GetConversation())

					channel.Broadcast.Send(interaction)
					// res, err := channel.conn.SendMessage(context.Background(), v.Info.Sender, &waProto.Message{
					// 	Conversation: proto.String("Hello, World!"),
					// })

					// fmt.Println(res)
					// fmt.Println(err)
				}
			}
		})

		// fmt.Println(channel.conn.Store.ID)
		if channel.conn.Store.ID == nil {
			// No ID stored, new login
			qrChan, _ := channel.conn.GetQRChannel(context.Background())
			err = channel.conn.Connect()
			if err != nil {
				fmt.Println(err)
			}
			for evt := range qrChan {
				if evt.Event == "code" {
					// Render the QR code here
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
					fmt.Println("QR code:", evt.Code)
				} else {
					fmt.Println("Login event:", evt.Event)
				}
			}
		} else {
			// Already logged in, just connect
			err = channel.conn.Connect()
			if err != nil {
				fmt.Println(err)
			}

			// channel.conn.handler
		}

		// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
		sc := make(chan os.Signal, 1)
		<-sc
	}()

	channel.running = true

	return nil
}

func (channel *WhatsappChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	channel.Broadcast.Close()
	channel.running = false

	return nil
}

func (channel *WhatsappChannel) SendAudio(*Interaction) error {
	panic("unimplemented")
}

func (channel *WhatsappChannel) SendButton(*Interaction) error {
	panic("unimplemented")
}

func (channel *WhatsappChannel) SendDocument(*Interaction) error {
	panic("unimplemented")
}

func (channel *WhatsappChannel) SendImage(*Interaction) error {
	panic("unimplemented")
}

func (channel *WhatsappChannel) SendText(interaction *Interaction) error {
	JID := interaction.Destination.Custom["JID"].(types.JID)
	_, err := channel.conn.SendMessage(context.Background(), JID, &waE2E.Message{
		Conversation: &interaction.Parameters.Text,
	})
	return err
}

func (channel *WhatsappChannel) SendVideo(*Interaction) error {
	panic("unimplemented")
}
