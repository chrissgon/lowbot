package lowbot

type Persist interface {
	Set(flow *Flow) error
	Get(sessionID string) (*Flow, error)
}
