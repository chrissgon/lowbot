package lowbot

import "fmt"

var (
	DEBUG                           = false
	ERR_CHANNEL_RUNNING             = fmt.Errorf("channel is running")
	ERR_CHANNEL_NOT_RUNNING         = fmt.Errorf("channel is not running")
	ERR_UNKNOWN_TELEGRAM_TOKEN      = fmt.Errorf("unknown telegram token")
	ERR_UNKNOWN_DISCORD_TOKEN       = fmt.Errorf("unknown discord token")
	ERR_UNKNOWN_CHATGPT_TOKEN       = fmt.Errorf("unknown chatgpt token")
	ERR_CONNECT_CHATGPT             = fmt.Errorf("connect to chatgpt failed")
	ERR_UNDEFINED_CHATGPT_ASSISTANT = fmt.Errorf("undefined chatgpt assistant")
	ERR_UNKNOWN_ACTION              = fmt.Errorf("unknown action")
	ERR_UNKNOWN_ROOM                = fmt.Errorf("unknown room")
	ERR_NIL_FLOW                    = fmt.Errorf("nil flow")
	ERR_NIL_STEP                    = fmt.Errorf("nil step")
	ERR_NIL_CHANNEL                 = fmt.Errorf("nil channel")
	ERR_UNKNOWN_DEFAULT_STEP        = fmt.Errorf("unknown step: default")
	ERR_UNKNOWN_INIT_STEP           = fmt.Errorf("unknown step: init")
	ERR_UNKNOWN_NEXT_STEP           = fmt.Errorf("unknown next step")
	ERR_INVALID_STEP                = fmt.Errorf("step invalid")
	ERR_PATTERN_NEXT_STEP           = fmt.Errorf("step pattern invalid")
	ERR_ENDED_FLOW                  = fmt.Errorf("flow ended")
	ERR_ROOM_STOPPED_FLOW           = fmt.Errorf("flow finished by room")
	ERR_FILE_NOT_PUBLIC             = fmt.Errorf("file is not public")

	ERR_FEATURE_UNIMPLEMENTED = fmt.Errorf("feature unimplemented")
)

const (
	CHANNEL_TELEGRAM_NAME        = "telegram"
	CHANNEL_WHATSAPP_TWILIO_NAME = "whatsapp-twilio"
	CHANNEL_WHATSAPP_DEVICE_NAME = "whatsapp-meow"
	CHANNEL_DISCORD_NAME         = "discord"
	FLOW_INIT_STEP_NAME          = "init"
	FLOW_END_STEP_NAME           = "end"
	FLOW_DEFAULT_STEP_NAME       = "default"
	FLOW_ERROR_STEP_NAME         = "error"
)
