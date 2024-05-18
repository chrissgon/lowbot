package lowbot

import (
	"os"
	"reflect"
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
	Audio    string   `yaml:"audio" json:"audio"`
	Buttons  []string `yaml:"buttons" json:"buttons"`
	Document string   `yaml:"document" json:"document"`
	Image    string   `yaml:"image" json:"image"`
	Text     string   `yaml:"text" json:"text"`
	Texts    []string `yaml:"texts"`
	Video    string   `yaml:"video" json:"video"`
}

type Steps map[string]*Step

func NewFlow(path string) (*Flow, error) {
	bytes, err := os.ReadFile(path)
	
	if err != nil {
		return nil, err
	}
	
	flow := &Flow{}

	err = yaml.Unmarshal(bytes, flow)

	return flow, err
}

func (flow *Flow) Start() {
	flow.Current = flow.Steps["init"]
}

func (flow *Flow) Next(in *Interaction) *Flow {
	if flow.NoHasNext() {
		return nil
	}

	for pattern, next := range flow.Current.Next {
		matched, _ := regexp.MatchString(pattern, in.Parameters.Text)

		if matched {
			flow.Current = flow.Steps[next]
			flow.Current.AddResponse(in)
			return flow
		}
	}

	flow.Current = flow.Steps[flow.Current.Next["default"]]
	flow.Current.AddResponse(in)
	return flow
}

func (flow *Flow) NoHasNext() bool {
	return flow.Current.Next == nil
}

func (flow *Flow) End() *Flow {
	flow.Current = flow.Steps["end"]
	return flow
}

func (flow *Flow) IsEnd() bool {
	endStep := flow.Steps["end"]

	return reflect.DeepEqual(endStep, flow.Current)
}

func (step *Step) AddResponse(in *Interaction) {
	step.Responses = append(step.Responses, in)
}

func (step *Step) GetLastResponse() *Interaction {
	return step.Responses[len(step.Responses)-1]
}

func (step *Step) GetLastResponseText() string {
	return step.Responses[len(step.Responses)-1].Parameters.Text
}
