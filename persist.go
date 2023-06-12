package lowbot

type Persist interface {
	Set(sessionID string, flow *Flow) error
	Get(sessionID string) (*Flow, error)
	Load() error
}
