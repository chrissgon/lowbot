package lowbot

type Who struct {
	WhoID  any
	Name   string
	Custom map[string]any
}

func NewWho(whoID any, name string) *Who {
	return &Who{
		WhoID:  whoID,
		Name:   name,
		Custom: map[string]any{},
	}
}
