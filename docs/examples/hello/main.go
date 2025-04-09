package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/chrissgon/lowbot"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func main() {
	lowbot.DEBUG = true

	// set custom actions
	lowbot.SetCustomActions(lowbot.ActionsMap{
		"TextUsername": func(interaction *lowbot.Interaction) (*lowbot.Interaction, bool) {
			template := lowbot.ParseTemplate(interaction.Step.Parameters.Texts)
			templateWithUsername := fmt.Sprintf(template, interaction.Parameters.Text)
			in := lowbot.NewInteractionMessageText(templateWithUsername)
			in.SetFrom(interaction.From)
			in.SetTo(interaction.To)
			return in, true
		},
	})

	// make a flow
	flow, _ := lowbot.NewFlow("./flow.yaml")

	fmt.Println()
	jid := types.NewJID("353892581511:5", "s.whatsapp.net")

	// fmt.Println(jid.IsBot())

	// make a channel. In this exemple is Telegram
	// channel, _ := lowbot.NewWhatsappTwilioChannel(os.Getenv("WHATSAPP_TWILIO_TOKEN"), os.Getenv("WHATSAPP_TWILIO_SID"))
	channel, err := lowbot.NewWhatsappDeviceChannel(&jid, func(evt whatsmeow.QRChannelItem, JID *types.JID) error {
		if evt.Event == "code" {
			png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)

			if err != nil {
				return err
			}

			base64Image := base64.StdEncoding.EncodeToString(png)

			fmt.Println(`{"qr":"data:image/png;base64,` + base64Image + `"}`)
			return nil
		}

		fmt.Println("JID", JID)

		return nil
	})

	if err != nil {
		panic(err)
	}
	// channel, _ := lowbot.NewWhatsappTwilioChannel(os.Getenv("WHATSAPP_TWILIO_TOKEN"), os.Getenv("WHATSAPP_TWILIO_SID"))

	// make a persist
	persist, _ := lowbot.NewMemoryFlowPersist()

	// make consumer
	consumer := lowbot.NewJourneyConsumer(flow, persist)

	// make bot
	bot := lowbot.NewBot(consumer, map[uuid.UUID]lowbot.IChannel{
		channel.GetChannel().ChannelID: channel,
	})

	// start bot
	bot.Start()

	// keep the process running
	sc := make(chan os.Signal, 1)
	<-sc
}
