package lowbot

import "errors"

var (
	ERR_UNKNOWN_TELEGRAM_TOKEN = errors.New("unknown TELEGRAM_TOKEN")
	ERR_UNKNOWN_DISCORD_TOKEN = errors.New("unknown DISCORD_TOKEN")
)