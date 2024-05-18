package lowbot

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Local struct {
	Sessions map[string]*Flow
}

var EnableLocalPersist = true

func (loc *Local) Load() error {
	file, err := os.Open("./local.json")

	if err != nil {
		return err
	}

	defer file.Close()

	bytes, _ := io.ReadAll(file)

	return json.Unmarshal(bytes, loc)
}

func (loc *Local) Get(sessionID string) (*Flow, error) {
	flow := loc.Sessions[sessionID]

	if flow == nil {
		return nil, fmt.Errorf("not found flow")
	}

	return flow, nil
}

func (loc *Local) Set(flow *Flow) error {
	loc.Sessions[flow.SessionID] = flow

	if EnableLocalPersist {
		go func() {
			file, _ := json.MarshalIndent(loc, "", " ")
			os.WriteFile("./local.json", file, 0644)
		}()
	}

	return nil
}

func NewLocalPersist() (Persist, error) {
	loc := &Local{Sessions: map[string]*Flow{}}

	if EnableLocalPersist {
		return loc, loc.Load()
	}

	return loc, nil
}
