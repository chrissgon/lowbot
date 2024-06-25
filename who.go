package lowbot

type Who struct {
	WhoID  string
	Name   string
	Custom map[string]any
}

func NewWho(whoID string, name string) *Who {
	return &Who{
		WhoID:  whoID,
		Name:   name,
		Custom: map[string]any{},
	}
}
