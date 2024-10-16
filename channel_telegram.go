package lowbot

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type TelegramChannel struct {
	*Channel
	conn    *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	running  bool
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

	channel := &TelegramChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_TELEGRAM_NAME,
			Broadcast: NewBroadcast[*Interaction](),
		},
		conn:   conn,
		running: false,
	}

	return channel, nil
}

func (channel *TelegramChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *TelegramChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	if channel.updates == nil {
		channel.updates = channel.conn.GetUpdatesChan(u)
	}

	go func() {
		for update := range channel.updates {
			if !channel.running {
				return
			}

			var interaction *Interaction

			if update.Message != nil {
				destination := NewWho(strconv.Itoa(int(update.Message.Chat.ID)), update.Message.From.UserName)
				interaction = NewInteractionMessageText(channel, destination, destination, update.Message.Text)
			}

			if update.CallbackQuery != nil {
				destination := NewWho(strconv.Itoa(int(update.CallbackQuery.From.ID)), update.CallbackQuery.From.UserName)
				interaction = NewInteractionMessageText(channel, destination, destination, update.CallbackData())
			}

			channel.Broadcast.Send(interaction)
		}
	}()

	channel.running = true

	return nil
}

func (channel *TelegramChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	channel.Broadcast.Close()
	channel.running = false
	return nil
}

func (channel *TelegramChannel) SendAudio(interaction *Interaction) error {
	channel.SendText(interaction)

	file := channel.getRequestFileDate(interaction.Parameters.File.GetFile().Path)
	chatID, err := strconv.Atoi(interaction.Destination.WhoID)

	if err != nil {
		return err
	}

	message := tgbotapi.NewAudio(int64(chatID), file)
	_, err = channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendButton(interaction *Interaction) error {
	button := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			channel.getButtons(interaction)...,
		),
	)

	chatID, err := strconv.Atoi(interaction.Destination.WhoID)

	if err != nil {
		return err
	}

	message := tgbotapi.NewMessage(int64(chatID), interaction.Parameters.Text)
	message.ReplyMarkup = button
	_, err = channel.conn.Send(message)

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
	chatID, err := strconv.Atoi(interaction.Destination.WhoID)

	if err != nil {
		return err
	}

	message := tgbotapi.NewDocument(int64(chatID), file)
	_, err = channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendImage(interaction *Interaction) error {
	channel.SendText(interaction)

	file := channel.getRequestFileDate(interaction.Parameters.File.GetFile().Path)
	chatID, err := strconv.Atoi(interaction.Destination.WhoID)

	if err != nil {
		return err
	}

	message := tgbotapi.NewPhoto(int64(chatID), file)
	_, err = channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendText(interaction *Interaction) error {
	chatID, err := strconv.Atoi(interaction.Destination.WhoID)

	if err != nil {
		return err
	}

	message := tgbotapi.NewMessage(int64(chatID), interaction.Parameters.Text)
	_, err = channel.conn.Send(message)

	return err
}

func (channel *TelegramChannel) SendVideo(interaction *Interaction) error {
	channel.SendText(interaction)

	file := channel.getRequestFileDate(interaction.Parameters.File.GetFile().Path)
	chatID, err := strconv.Atoi(interaction.Destination.WhoID)

	if err != nil {
		return err
	}

	message := tgbotapi.NewDocument(int64(chatID), file)
	_, err = channel.conn.Send(message)

	return err
}

func (*TelegramChannel) getRequestFileDate(str string) (file tgbotapi.RequestFileData) {
	if IsURL(str) {
		file = tgbotapi.FileURL(str)
	}

	file = tgbotapi.FilePath(str)

	return
}
