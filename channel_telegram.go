package lowbot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type TelegramChannel struct {
	*Channel
	conn *tgbotapi.BotAPI
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
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_TELEGRAM_NAME,
		},
		conn: conn,
	}, nil
}

func (channel *TelegramChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *TelegramChannel) Next(interaction chan *Interaction) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := channel.conn.GetUpdatesChan(u)

	for update := range updates {
		var i *Interaction

		if update.Message != nil {
			sender := NewWho(update.Message.Chat.ID, update.Message.From.UserName)
			i = NewInteractionMessageText(channel, sender, update.Message.Text)
		}

		if update.CallbackQuery != nil {
			sender := NewWho(update.CallbackQuery.From.ID, update.CallbackQuery.From.UserName)
			i = NewInteractionMessageText(channel, sender, update.CallbackData())
		}

		interaction <- i
	}
}

func (channel *TelegramChannel) SendAudio(interaction *Interaction) error {
	channel.SendText(interaction)

	file := channel.getRequestFileDate(interaction.Parameters.File.GetFile().Path)
	chatID := interaction.Sender.WhoID.(int64)
	message := tgbotapi.NewAudio(chatID, file)
	_, err := channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendButton(interaction *Interaction) error {
	button := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			channel.getButtons(interaction)...,
		),
	)

	chatID := interaction.Sender.WhoID.(int64)
	message := tgbotapi.NewMessage(chatID, interaction.Parameters.Text)
	message.ReplyMarkup = button
	_, err := channel.conn.Send(message)

	return err
}

func (*TelegramChannel) getButtons(interaction *Interaction) (buttons []tgbotapi.InlineKeyboardButton) {
	for _, button := range interaction.Parameters.Buttons {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(button, button))
	}
	return
}

func (channel *TelegramChannel) SendDocument(interaction *Interaction) error {
	channel.SendText(interaction)

	file := channel.getRequestFileDate(interaction.Parameters.File.GetFile().Path)
	chatID := interaction.Sender.WhoID.(int64)
	message := tgbotapi.NewDocument(chatID, file)
	_, err := channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendImage(interaction *Interaction) error {
	channel.SendText(interaction)

	file := channel.getRequestFileDate(interaction.Parameters.File.GetFile().Path)
	chatID := interaction.Sender.WhoID.(int64)
	message := tgbotapi.NewPhoto(chatID, file)
	_, err := channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendText(interaction *Interaction) error {
	chatID := interaction.Sender.WhoID.(int64)
	message := tgbotapi.NewMessage(chatID, interaction.Parameters.Text)
	_, err := channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendVideo(interaction *Interaction) error {
	channel.SendText(interaction)

	file := channel.getRequestFileDate(interaction.Parameters.File.GetFile().Path)
	chatID := interaction.Sender.WhoID.(int64)
	message := tgbotapi.NewDocument(chatID, file)
	_, err := channel.conn.Send(message)

	return err
}

func (*TelegramChannel) getRequestFileDate(str string) (file tgbotapi.RequestFileData) {
	if IsURL(str) {
		file = tgbotapi.FileURL(str)
	}

	file = tgbotapi.FilePath(str)

	return
}
