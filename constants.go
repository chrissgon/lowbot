package lowbot

import "errors"

var (
	DEBUG                           = false
	ERR_CHANNEL_RUNNING             = errors.New("channel is running")
	ERR_CHANNEL_NOT_RUNNING         = errors.New("channel is not running")
	ERR_UNKNOWN_TELEGRAM_TOKEN      = errors.New("unknown telegram token")
	ERR_UNKNOWN_DISCORD_TOKEN       = errors.New("unknown discord token")
	ERR_UNKNOWN_CHATGPT_TOKEN       = errors.New("unknown chatgpt token")
	ERR_CONNECT_CHATGPT             = errors.New("connect to chatgpt failed")
	ERR_UNDEFINED_CHATGPT_ASSISTANT = errors.New("undefined chatgpt assistant")
	ERR_UNKNOWN_ACTION              = errors.New("unknown action")
	ERR_UNKNOWN_ROOM                = errors.New("unknown room")
	ERR_NIL_FLOW                    = errors.New("nil flow")
	ERR_NIL_STEP                    = errors.New("nil step")
	ERR_NIL_CHANNEL                 = errors.New("nil channel")
	ERR_UNKNOWN_DEFAULT_STEP        = errors.New("unknown step: default")
	ERR_UNKNOWN_INIT_STEP           = errors.New("unknown step: init")
	ERR_UNKNOWN_NEXT_STEP           = errors.New("unknown next step")
	ERR_PATTERN_NEXT_STEP           = errors.New("step pattern invalid")
	ERR_ENDED_FLOW                  = errors.New("flow ended")
	ERR_ROOM_STOPPED_FLOW           = errors.New("flow finished by room")
	ERR_FILE_NOT_PUBLIC           = errors.New("file is not public")

	ERR_FEATURE_UNIMPLEMENTED = errors.New("feature unimplemented")
)

const (
	CHANNEL_TELEGRAM_NAME        = "telegram"
	CHANNEL_WHATSAPP_TWILIO_NAME = "whatsapp twilio"
	CHANNEL_DISCORD_NAME         = "discord"
	CONSUMER_JOURNEY_NAME        = "journey"
	CONSUMER_CHATGPT_NAME        = "chatgpt"
	FLOW_INIT_STEP_NAME          = "init"
	FLOW_END_STEP_NAME           = "end"
	FLOW_DEFAULT_STEP_NAME       = "default"
)
