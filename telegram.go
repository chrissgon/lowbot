package lowbot

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	conn *tgbotapi.BotAPI
}

func (tg *Telegram) SendAudio(in Interaction) error {
	tg.SendText(in)

	file := tg.getRequestFileDate(in.Parameters.Audio)

	_, err := tg.conn.Send(tgbotapi.NewAudio(StringToInt64(in.SessionID), file))

	return err
}

func (tg *Telegram) SendButton(in Interaction) error {
	button := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tg.getButtons(in)...,
		),
	)

	message := tgbotapi.NewMessage(StringToInt64(in.SessionID), in.Parameters.Text)
	message.ReplyMarkup = button

	_, err := tg.conn.Send(message)
	return err
}

func (*Telegram) getButtons(in Interaction) (buttons []tgbotapi.InlineKeyboardButton) {
	for _, button := range in.Parameters.Buttons {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	return
}

func (tg *Telegram) SendDocument(in Interaction) error {
	tg.SendText(in)

	file := tg.getRequestFileDate(in.Parameters.Document)

	_, err := tg.conn.Send(tgbotapi.NewDocument(StringToInt64(in.SessionID), file))

	return err
}

func (tg *Telegram) SendImage(in Interaction) error {
	tg.SendText(in)

	file := tg.getRequestFileDate(in.Parameters.Image)

	_, err := tg.conn.Send(tgbotapi.NewPhoto(StringToInt64(in.SessionID), file))

	return err
}

func (tg *Telegram) SendText(in Interaction) error {
	_, err := tg.conn.Send(tgbotapi.NewMessage(StringToInt64(in.SessionID), in.Parameters.Text))
	return err
}

func (tg *Telegram) SendVideo(in Interaction) error {
	tg.SendText(in)

	file := tg.getRequestFileDate(in.Parameters.Video)

	_, err := tg.conn.Send(tgbotapi.NewVideo(StringToInt64(in.SessionID), file))

	return err
}

func (tg *Telegram) Next(in chan Interaction) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tg.conn.GetUpdatesChan(u)

	for update := range updates {
		var i Interaction

		if update.Message != nil {
			i = NewInteractionMessageText(Int64ToString(update.Message.Chat.ID), update.Message.Text)
		}

		if update.CallbackQuery != nil {
			i = NewInteractionMessageText(Int64ToString(update.CallbackQuery.From.ID), update.CallbackData())
		}

		in <- i
	}
}

func (*Telegram) getRequestFileDate(str string) (file tgbotapi.RequestFileData) {
	file = tgbotapi.FilePath(str)

	if IsURL(str) {
		file = tgbotapi.FileURL(str)
	}
	
	return
}

func NewTelegram() (Channel, error) {
	token := os.Getenv("TELEGRAM_TOKEN")

	if token == "" {
		return nil, NewError("NewTelegram", ERR_UNKNOWN_TELEGRAM_TOKEN)
	}

	conn, err := tgbotapi.NewBotAPI(token)
	conn.Debug = false

	if err != nil {
		return nil, NewError("NewTelegram", err)
	}

	return &Telegram{conn: conn}, nil
}
