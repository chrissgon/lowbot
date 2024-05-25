package lowbot

import "errors"

var (
	DEBUG                        = true
	ERR_UNKNOWN_INTERACTION_TYPE = errors.New("unknown interaction type")
	ERR_UNKNOWN_TELEGRAM_TOKEN   = errors.New("unknown telegram token")
	ERR_UNKNOWN_DISCORD_TOKEN    = errors.New("unknown discord token")
	ERR_UNKNOWN_CHATGPT_TOKEN    = errors.New("unknown chatgpt token")
	ERR_CONNECT_CHATGPT          = errors.New("connect to chatgpt failed")
	ERR_UNKNOWN_ACTION           = errors.New("unknown action")
	ERR_NIL_FLOW                 = errors.New("nil flow")
	ERR_NIL_STEP                 = errors.New("nil step")
	ERR_UNKNOWN_DEFAULT_STEP     = errors.New("unknown step: default")
	ERR_UNKNOWN_INIT_STEP        = errors.New("unknown step: init")
	ERR_UNKNOWN_NEXT_STEP        = errors.New("unknown next step")
	ERR_PATTERN_NEXT_STEP        = errors.New("step pattern invalid")
)
