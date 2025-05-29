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
	Replier         string `yaml:"replier" json:"replier"`
	Description     string `yaml:"description" json:"description"`
	Steps           Steps  `yaml:"steps" json:"steps"`
	CurrentStep     Step
	CurrentStepName string

	Waiting bool
}

type Step struct {
	Action     string            `yaml:"action" json:"action"`
	Next       map[string]string `yaml:"next" json:"next"`
	Parameters StepParameters    `yaml:"parameters" json:"parameters"`
}

type StepParameters struct {
	Buttons []string          `yaml:"buttons" json:"buttons"`
	Path    string            `yaml:"path" json:"path"`
	URL     string            `yaml:"url" json:"url"`
	Text    string            `yaml:"text" json:"text"`
	Texts   []string          `yaml:"texts" json:"texts"`
	Headers map[string]string `yaml:"headers" json:"headers"`
	Timeout int               `yaml:"timeout" json:"timeout"`
	Custom  map[string]any    `yaml:"custom" json:"custom"`
}

type Steps map[string]Step

type IFlowPersist interface {
	Set(any, *Flow) error
	Get(any) (*Flow, error)
}

var FlowPersist = NewMemoryFlowPersist()

func GetCurrentStep(sessionID any) (Step, error) {
	flow, err := FlowPersist.Get(sessionID)

	if err != nil {
		return Step{}, err
	}

	if flow == nil {
		return Step{}, ERR_NIL_FLOW
	}

	return flow.CurrentStep, nil
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
	// flow.FlowID = uuid.New()
	step, exists := flow.Steps[FLOW_INIT_STEP_NAME]

	if !exists {
		return ERR_UNKNOWN_INIT_STEP
	}

	flow.CurrentStep = step
	flow.CurrentStepName = FLOW_INIT_STEP_NAME

	return nil
}

func (flow *Flow) NextError() {
	next := flow.CurrentStep.Next[FLOW_ERROR_STEP_NAME]
	step, exists := flow.Steps[next]

	if !exists {
		return
	}

	flow.CurrentStep = step
	flow.CurrentStepName = next

	return
}

func (flow *Flow) Next(interaction Interaction) error {
	if flow.Ended() {
		return ERR_ENDED_FLOW
	}

	err := flow.goNextStep(interaction)

	if err != nil {
		return err
	}

	return nil
}

func (flow *Flow) goNextStep(interaction Interaction) error {
	if flow.NoHasNext() {
		return ERR_UNKNOWN_NEXT_STEP
	}

	for pattern, next := range flow.CurrentStep.Next {
		matched, err := regexp.MatchString(pattern, interaction.Parameters.Text)

		if err != nil {
			return ERR_PATTERN_NEXT_STEP
		}

		step, exists := flow.Steps[next]

		if matched && !exists {
			return ERR_INVALID_STEP
		}

		if matched {
			flow.CurrentStep = step
			flow.CurrentStepName = next

			return nil
		}
	}

	next, exists := flow.CurrentStep.Next[FLOW_DEFAULT_STEP_NAME]

	if !exists {
		return ERR_UNKNOWN_DEFAULT_STEP
	}

	step := flow.Steps[next]

	flow.CurrentStep = step
	flow.CurrentStepName = next

	return nil
}

func (flow *Flow) NoHasNext() bool {
	return flow.CurrentStep.Next == nil
}

func (flow *Flow) Ended() bool {
	return flow.CurrentStepName == FLOW_END_STEP_NAME
}

func (flow *Flow) Wait(interaction Interaction) error {
	flow.Waiting = true
	return FlowPersist.Set(interaction.From.WhoID, flow)
}

func (flow *Flow) Continue(interaction Interaction) error {
	flow.Waiting = false
	return FlowPersist.Set(interaction.From.WhoID, flow)
}
