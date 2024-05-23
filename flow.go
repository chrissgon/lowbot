package lowbot

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Flow struct {
	SessionID string
	Name      string `yaml:"name"`
	Steps     Steps  `yaml:"steps"`
	Current   *Step
}

type Step struct {
	Action     string            `yaml:"action"`
	Next       map[string]string `yaml:"next"`
	Parameters StepParameters    `yaml:"parameters"`
	Responses  []*Interaction
}

type StepParameters struct {
	Audio    string         `yaml:"audio"`
	Buttons  []string       `yaml:"buttons"`
	Document string         `yaml:"document"`
	Image    string         `yaml:"image"`
	Text     string         `yaml:"text"`
	Texts    []string       `yaml:"texts"`
	Video    string         `yaml:"video"`
	Custom   map[string]any `yaml:"custom"`
}

type Steps map[string]*Step

type FlowPersist interface {
	Set(flow *Flow) error
	Get(sessionID string) (*Flow, error)
}

func NewFlow(path string) (*Flow, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	flow := &Flow{}

	err = yaml.Unmarshal(bytes, flow)

	return flow, err
}

func (flow *Flow) Start() error {
	step, exists := flow.Steps["init"]

	if !exists {
		return ERR_UNKNOWN_INIT_STEP
	}

	flow.Current = step

	return nil
}

func (flow *Flow) Next(interaction *Interaction) error {
	if flow.NoHasNext() {
		return ERR_UNKNOWN_NEXT_STEP
	}

	err := flow.setNextFlow(interaction)

	if err != nil {
		return err
	}

	flow.Current.AddResponse(interaction)

	return nil
}

func (flow *Flow) setNextFlow(interaction *Interaction) error {
	for pattern, next := range flow.Current.Next {
		matched, err := regexp.MatchString(pattern, interaction.Parameters.Text)

		if err != nil {
			return ERR_PATTERN_NEXT_STEP
		}

		if matched {
			flow.Current = flow.Steps[next]

			return nil
		}
	}

	next, exists := flow.Current.Next["default"]

	if !exists {
		return ERR_UNKNOWN_DEFAULT_STEP
	}

	flow.Current = flow.Steps[next]

	return nil
}

func (flow *Flow) NoHasNext() bool {
	return flow.Current.Next == nil
}

func (step *Step) AddResponse(interaction *Interaction) {
	step.Responses = append(step.Responses, interaction)
}

func (step *Step) GetLastResponse() *Interaction {
	return step.Responses[len(step.Responses)-1]
}

func (step *Step) GetLastResponseText() string {
	return step.Responses[len(step.Responses)-1].Parameters.Text
}
