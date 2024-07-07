package lowbot

type Guest struct {
	Who     *Who
	Channel IChannel
}

func NewGuest(who *Who, channel IChannel) *Guest {
	return &Guest{
		Who:     who,
		Channel: channel,
	}
}
