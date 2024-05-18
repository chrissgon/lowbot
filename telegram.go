package lowbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	conn *tgbotapi.BotAPI
}

func NewTelegram(token string) (Channel, error) {
	if token == "" {
		return nil, ERR_UNKNOWN_TELEGRAM_TOKEN
	}

	conn, err := tgbotapi.NewBotAPI(token)
	conn.Debug = false

	if err != nil {
		return nil, err
	}

	return &Telegram{conn: conn}, nil
}

func (tg *Telegram) SendAudio(interaction *Interaction) error {
	tg.SendText(interaction)

	file := tg.getRequestFileDate(interaction.Parameters.Audio)

	_, err := tg.conn.Send(tgbotapi.NewAudio(StringToInt64(interaction.SessionID), file))

	return err
}

func (tg *Telegram) SendButton(interaction *Interaction) error {
	button := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tg.getButtons(interaction)...,
		),
	)

	message := tgbotapi.NewMessage(StringToInt64(interaction.SessionID), interaction.Parameters.Text)
	message.ReplyMarkup = button

	_, err := tg.conn.Send(message)
	return err
}

func (*Telegram) getButtons(interaction *Interaction) (buttons []tgbotapi.InlineKeyboardButton) {
	for _, button := range interaction.Parameters.Buttons {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	return
}

func (tg *Telegram) SendDocument(interaction *Interaction) error {
	tg.SendText(interaction)

	file := tg.getRequestFileDate(interaction.Parameters.Document)

	_, err := tg.conn.Send(tgbotapi.NewDocument(StringToInt64(interaction.SessionID), file))

	return err
}

func (tg *Telegram) SendImage(interaction *Interaction) error {
	tg.SendText(interaction)

	file := tg.getRequestFileDate(interaction.Parameters.Image)

	_, err := tg.conn.Send(tgbotapi.NewPhoto(StringToInt64(interaction.SessionID), file))

	return err
}

func (tg *Telegram) SendText(interaction *Interaction) error {
	_, err := tg.conn.Send(tgbotapi.NewMessage(StringToInt64(interaction.SessionID), interaction.Parameters.Text))
	return err
}

func (tg *Telegram) SendVideo(interaction *Interaction) error {
	tg.SendText(interaction)

	file := tg.getRequestFileDate(interaction.Parameters.Video)

	_, err := tg.conn.Send(tgbotapi.NewVideo(StringToInt64(interaction.SessionID), file))

	return err
}

func (tg *Telegram) Next(interaction chan *Interaction) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tg.conn.GetUpdatesChan(u)

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

func (*Telegram) getRequestFileDate(str string) (file tgbotapi.RequestFileData) {
	file = tgbotapi.FilePath(str)

	if IsURL(str) {
		file = tgbotapi.FileURL(str)
	}

	return
}
