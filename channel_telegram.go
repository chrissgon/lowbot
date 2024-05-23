package lowbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type TelegramChannel struct {
	channelID uuid.UUID
	conn      *tgbotapi.BotAPI
}

func NewTelegramChannel(token string) (IChannel, error) {
	if token == "" {
		return nil, ERR_UNKNOWN_TELEGRAM_TOKEN
	}

	conn, err := tgbotapi.NewBotAPI(token)
	conn.Debug = false

	if err != nil {
		return nil, err
	}

	return &TelegramChannel{
		channelID: uuid.New(),
		conn:      conn,
	}, nil
}

func (telegram *TelegramChannel) ChannelID() uuid.UUID {
	return telegram.channelID
}

func (telegram *TelegramChannel) SendAudio(interaction *Interaction) error {
	telegram.SendText(interaction)

	file := telegram.getRequestFileDate(interaction.Parameters.Audio)

	_, err := telegram.conn.Send(tgbotapi.NewAudio(StringToInt64(interaction.SessionID), file))

	return err
}

func (telegram *TelegramChannel) SendButton(interaction *Interaction) error {
	button := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			telegram.getButtons(interaction)...,
		),
	)

	message := tgbotapi.NewMessage(StringToInt64(interaction.SessionID), interaction.Parameters.Text)
	message.ReplyMarkup = button

	_, err := telegram.conn.Send(message)
	return err
}

func (*TelegramChannel) getButtons(interaction *Interaction) (buttons []tgbotapi.InlineKeyboardButton) {
	for _, button := range interaction.Parameters.Buttons {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	return
}

func (telegram *TelegramChannel) SendDocument(interaction *Interaction) error {
	telegram.SendText(interaction)

	file := telegram.getRequestFileDate(interaction.Parameters.Document)

	_, err := telegram.conn.Send(tgbotapi.NewDocument(StringToInt64(interaction.SessionID), file))

	return err
}

func (telegram *TelegramChannel) SendImage(interaction *Interaction) error {
	telegram.SendText(interaction)

	file := telegram.getRequestFileDate(interaction.Parameters.Image)

	_, err := telegram.conn.Send(tgbotapi.NewPhoto(StringToInt64(interaction.SessionID), file))

	return err
}

func (telegram *TelegramChannel) SendText(interaction *Interaction) error {
	_, err := telegram.conn.Send(tgbotapi.NewMessage(StringToInt64(interaction.SessionID), interaction.Parameters.Text))
	return err
}

func (telegram *TelegramChannel) SendVideo(interaction *Interaction) error {
	telegram.SendText(interaction)

	file := telegram.getRequestFileDate(interaction.Parameters.Video)

	_, err := telegram.conn.Send(tgbotapi.NewVideo(StringToInt64(interaction.SessionID), file))

	return err
}

func (telegram *TelegramChannel) Next(interaction chan *Interaction) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := telegram.conn.GetUpdatesChan(u)

	for update := range updates {
		var i *Interaction

		if update.Message != nil {
			i = NewInteractionMessageText(Int64ToString(update.Message.Chat.ID), update.Message.Text)
		}

		if update.CallbackQuery != nil {
			i = NewInteractionMessageText(Int64ToString(update.CallbackQuery.From.ID), update.CallbackData())
		}

		interaction <- i
	}
}

func (*TelegramChannel) getRequestFileDate(str string) (file tgbotapi.RequestFileData) {
	file = tgbotapi.FilePath(str)

	if IsURL(str) {
		file = tgbotapi.FileURL(str)
	}

	return
}
