package lowbot

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Flow struct {
	FlowID          uuid.UUID
	Name            string `yaml:"name" json:"name"`
	Description     string `yaml:"description" json:"description"`
	Steps           Steps  `yaml:"steps" json:"steps"`
	CurrentStep     *Step
	CurrentStepName string
}

type Step struct {
	Action     string            `yaml:"action" json:"action"`
	Next       map[string]string `yaml:"next" json:"next"`
	Parameters StepParameters    `yaml:"parameters" json:"parameters"`
}

type StepParameters struct {
	Buttons []string       `yaml:"buttons" json:"buttons"`
	Path    string         `yaml:"path" json:"path"`
	Text    string         `yaml:"text" json:"text"`
	Texts   []string       `yaml:"texts" json:"texts"`
	Custom  map[string]any `yaml:"custom" json:"custom"`
}

type Steps map[string]*Step

type FlowPersist interface {
	Set(any, *Flow) error
	Get(any) (*Flow, error)
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

func NewFlowByJSON(strJSON string) (*Flow, error) {
	flow := &Flow{}

	err := json.Unmarshal([]byte(strJSON), flow)

	return flow, err
}

func NewFlowByJSONFile(path string) (*Flow, error) {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	flow := &Flow{}

	err = json.Unmarshal(bytes, flow)

	return flow, err
}

func (flow *Flow) Start() error {
	flow.FlowID = uuid.New()

	step, exists := flow.Steps[FLOW_INIT_STEP_NAME]

	if !exists {
		return ERR_UNKNOWN_INIT_STEP
	}

	flow.CurrentStep = step
	flow.CurrentStepName = FLOW_INIT_STEP_NAME

	return nil
}

func (flow *Flow) Next(interaction *Interaction) error {
	if flow.Ended() {
		return ERR_ENDED_FLOW
	}
	if flow.NoHasNext() {
		return ERR_UNKNOWN_NEXT_STEP
	}

	err := flow.goNextStep(interaction)

	if err != nil {
		return err
	}

	return nil
}

func (flow *Flow) goNextStep(interaction *Interaction) error {
	for pattern, next := range flow.CurrentStep.Next {
		matched, err := regexp.MatchString(pattern, interaction.Parameters.Text)

		if err != nil {
			return ERR_PATTERN_NEXT_STEP
		}

		if matched {
			flow.CurrentStep = flow.Steps[next]
			flow.CurrentStepName = next

			return nil
		}
	}

	next, exists := flow.CurrentStep.Next[FLOW_DEFAULT_STEP_NAME]

	if !exists {
		return ERR_UNKNOWN_DEFAULT_STEP
	}

	flow.CurrentStep = flow.Steps[next]
	flow.CurrentStepName = next

	return nil
}

func (flow *Flow) NoHasNext() bool {
	return flow.CurrentStep.Next == nil
}

func (flow *Flow) Ended() bool {
	return flow.CurrentStepName == FLOW_END_STEP_NAME
}
