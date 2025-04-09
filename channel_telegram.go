package lowbot

import (
	"bytes"
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

type TelegramChannel struct {
	*Channel
	running bool
	conn    *bot.Bot
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewTelegramChannel(token string) (IChannel, error) {
	if token == "" {
		return nil, ERR_UNKNOWN_TELEGRAM_TOKEN
	}

	channel := &TelegramChannel{
		Channel: &Channel{
			ChannelID: uuid.New(),
			Name:      CHANNEL_TELEGRAM_NAME,
			Broadcast: NewBroadcast[*Interaction](),
		},
		running: false,
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(channel.telegramMessageHandler),
		bot.WithCallbackQueryDataHandler("button", bot.MatchTypePrefix, channel.telegramCallbackHandler),
	}

	bot, err := bot.New(token, opts...)

	if err != nil {
		return nil, err
	}

	channel.conn = bot

	return channel, nil
}

func (channel *TelegramChannel) telegramMessageHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !channel.running {
		return
	}

	if update.Message != nil {
		from := NewWho(strconv.Itoa(int(update.Message.Chat.ID)), update.Message.From.Username)
		interaction := NewInteractionMessageText(update.Message.Text)
		interaction.SetFrom(from)
		channel.Broadcast.Send(interaction)
		return
	}

	from := NewWho(strconv.Itoa(int(update.CallbackQuery.From.ID)), update.CallbackQuery.From.Username)
	interaction := NewInteractionMessageText(update.CallbackQuery.Data)
	interaction.SetFrom(from)
	channel.Broadcast.Send(interaction)
}

func (channel *TelegramChannel) telegramCallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if !channel.running {
		return
	}

	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
}

func (channel *TelegramChannel) GetChannel() *Channel {
	return channel.Channel
}

func (channel *TelegramChannel) Start() error {
	if channel.running {
		return ERR_CHANNEL_RUNNING
	}

	channel.ctx, channel.cancel = context.WithCancel(context.Background())

	go channel.conn.Start(channel.ctx)

	channel.running = true

	return nil
}

func (channel *TelegramChannel) Stop() error {
	if !channel.running {
		return ERR_CHANNEL_NOT_RUNNING
	}

	channel.Broadcast.Close()
	channel.cancel()
	channel.running = false

	return nil
}

func (channel *TelegramChannel) SendAudio(interaction *Interaction) error {
	chatID, err := strconv.Atoi(interaction.From.WhoID)

	if err != nil {
		return err
	}

	err = interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	_, err = channel.conn.SendAudio(channel.ctx, &bot.SendAudioParams{
		ChatID:  int64(chatID),
		Caption: interaction.Parameters.Text,
		Audio: &models.InputFileUpload{
			Data:     bytes.NewReader(interaction.Parameters.File.GetFile().Bytes),
			Filename: interaction.Parameters.File.GetFile().Name,
		},
	})

	return err
}

func (channel *TelegramChannel) SendButton(interaction *Interaction) error {
	chatID, err := strconv.Atoi(interaction.From.WhoID)

	if err != nil {
		return err
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			channel.getButtons(interaction),
		},
	}

	_, err = channel.conn.SendMessage(channel.ctx, &bot.SendMessageParams{
		ChatID:      int64(chatID),
		Text:        escapeMarkdownV2(interaction.Parameters.Text),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: kb,
	})

	return err
}

func (*TelegramChannel) getButtons(interaction *Interaction) (buttons []models.InlineKeyboardButton) {
	for i, button := range interaction.Parameters.Buttons {
		buttons = append(buttons, models.InlineKeyboardButton{
			Text:         button,
			CallbackData: strconv.Itoa(i + 1),
		})
	}
	return
}

func (channel *TelegramChannel) SendDocument(interaction *Interaction) error {
	chatID, err := strconv.Atoi(interaction.From.WhoID)

	if err != nil {
		return err
	}

	err = interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	_, err = channel.conn.SendDocument(channel.ctx, &bot.SendDocumentParams{
		ChatID:  int64(chatID),
		Caption: interaction.Parameters.Text,
		Document: &models.InputFileUpload{
			Data:     bytes.NewReader(interaction.Parameters.File.GetFile().Bytes),
			Filename: interaction.Parameters.File.GetFile().Name,
		},
	})

	return err
}

func (channel *TelegramChannel) SendImage(interaction *Interaction) error {
	chatID, err := strconv.Atoi(interaction.From.WhoID)

	if err != nil {
		return err
	}

	err = interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	_, err = channel.conn.SendPhoto(channel.ctx, &bot.SendPhotoParams{
		ChatID:  int64(chatID),
		Caption: interaction.Parameters.Text,
		Photo: &models.InputFileUpload{
			Data:     bytes.NewReader(interaction.Parameters.File.GetFile().Bytes),
			Filename: interaction.Parameters.File.GetFile().Name,
		},
	})

	return err
}

func (channel *TelegramChannel) SendText(interaction *Interaction) error {
	chatID, err := strconv.Atoi(interaction.From.WhoID)

	if err != nil {
		return err
	}

	_, err = channel.conn.SendMessage(channel.ctx, &bot.SendMessageParams{
		ChatID:    int64(chatID),
		Text:      escapeMarkdownV2(interaction.Parameters.Text),
		ParseMode: models.ParseModeMarkdown,
	})

	return err
}

func (channel *TelegramChannel) SendVideo(interaction *Interaction) error {
	chatID, err := strconv.Atoi(interaction.From.WhoID)

	if err != nil {
		return err
	}

	err = interaction.Parameters.File.Read()

	if err != nil {
		return err
	}

	_, err = channel.conn.SendVideo(channel.ctx, &bot.SendVideoParams{
		ChatID:  int64(chatID),
		Caption: interaction.Parameters.Text,
		Video: &models.InputFileUpload{
			Data:     bytes.NewReader(interaction.Parameters.File.GetFile().Bytes),
			Filename: interaction.Parameters.File.GetFile().Name,
		},
	})

	return err
}

func escapeMarkdownV2(text string) string {
	linkRegex := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	links := linkRegex.FindAllStringSubmatchIndex(text, -1)

	if len(links) == 0 {
		return escapeText(text)
	}

	var result strings.Builder
	lastIndex := 0

	for _, match := range links {
		result.WriteString(escapeText(text[lastIndex:match[0]]))
		linkText := escapeText(text[match[2]:match[3]])
		linkURL := escapeURL(text[match[4]:match[5]])
		result.WriteString("[" + linkText + "](" + linkURL + ")")
		lastIndex = match[1]
	}

	if lastIndex < len(text) {
		result.WriteString(escapeText(text[lastIndex:]))
	}

	return result.String()
}

func escapeText(text string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\", "_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]",
		"(", "\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
		"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|", "\\|",
		"{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!")
	return replacer.Replace(text)
}

func escapeURL(url string) string {
	replacer := strings.NewReplacer("\\", "\\\\", ")", "\\)", ".", "\\.", "!", "\\!")
	return replacer.Replace(url)
}
