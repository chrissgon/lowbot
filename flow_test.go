package lowbot

import (
	"errors"
	"reflect"
	"testing"
)

func TestFlow_NewFlow(t *testing.T) {
	_, err := NewFlow("")

	if err == nil {
		t.Errorf(FormatTestError("any error", nil))
	}

	expect := newFlowMock()
	have, err := NewFlow("./mocks/flow.yaml")

	if err != nil {
		t.Errorf(FormatTestError(nil, err))
	}

	if !reflect.DeepEqual(expect, *have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFlow_StartFlow(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")

	delete(flow.Steps, "init")

	err := flow.Start()

	if !errors.Is(err, ERR_UNKNOWN_INIT_STEP) {
		t.Errorf(FormatTestError(ERR_UNKNOWN_INIT_STEP, err))
	}

	expect := flow.Steps["init"]
	have := flow.CurrentStep

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFlow_NextFlow(t *testing.T) {
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")

	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.CurrentStepName = FLOW_END_STEP_NAME
	err := flow.Next(interaction)

	if !errors.Is(err, ERR_ENDED_FLOW) {
		t.Errorf(FormatTestError(ERR_ENDED_FLOW, err))
	}

	flow, _ = NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.CurrentStep.Next = nil
	err = flow.Next(interaction)

	if !errors.Is(err, ERR_UNKNOWN_NEXT_STEP) {
		t.Errorf(FormatTestError(ERR_UNKNOWN_NEXT_STEP, err))
	}

	flow, _ = NewFlow("./mocks/flow.yaml")
	flow.Start()
	delete(flow.CurrentStep.Next, "default")
	err = flow.Next(interaction)

	if !errors.Is(err, ERR_UNKNOWN_DEFAULT_STEP) {
		t.Errorf(FormatTestError(ERR_UNKNOWN_DEFAULT_STEP, err))
	}

	flow, _ = NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, ""))

	expect := flow.Steps["audio"]
	have := flow.CurrentStep

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFlow_NoHasNext(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.CurrentStep.Next = nil

	if !flow.NoHasNext() {
		t.Errorf(FormatTestError(true, false))
	}
}

func TestFlow_AddResponse(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, ""))

	expect := 1
	have := len(flow.Responses)

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
func TestFlow_AddResponseValue(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")

	flow.CurrentStepName = FLOW_END_STEP_NAME

	expect := true
	have := flow.Ended()

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
func TestFlow_GetLastResponse(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()

	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow.Next(interaction)

	expect := interaction
	have := flow.GetLastResponse()

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}
func TestFlow_GetLastResponsetText(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()

	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "")
	flow.Next(interaction)

	expect := interaction.Parameters.Text
	have := flow.GetLastResponseText()

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
func TestFlow_Ended(t *testing.T) {
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "Response")
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(interaction)

	expect := interaction
	have := flow.Responses[0]

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func newFlowMock() Flow {
	return Flow{
		Name: "flow",
		Steps: Steps{
			"init": {
				Next: map[string]string{
					"default": "audio",
				},
			},
			"audio": {
				Next: map[string]string{
					"default": "button",
				},
				Action: "File",
				Parameters: StepParameters{
					Path: "./mocks/music.mp3",
				},
			},
			"button": {
				Next: map[string]string{
					"default": "document",
				},
				Action: "Button",
				Parameters: StepParameters{
					Buttons: []string{"yes", "no"},
					Texts:   []string{"buttons here"},
				},
			},
			"document": {
				Next: map[string]string{
					"default": "image",
				},
				Action: "File",
				Parameters: StepParameters{
					Path: "./mocks/features.txt",
				},
			},
			"image": {
				Next: map[string]string{
					"default": "text",
				},
				Action: "File",
				Parameters: StepParameters{
					Path: "./mocks/image.jpg",
				},
			},
			"text": {
				Next: map[string]string{
					"default": "video",
				},
				Action: "Text",
				Parameters: StepParameters{
					Texts: []string{"texts"},
				},
			},
			"video": {
				Next: map[string]string{
					"default": "end",
				},
				Action: "File",
				Parameters: StepParameters{
					Path: "./mocks/video.mp4",
				},
			},
			"end": {
				Action: "Text",
				Parameters: StepParameters{
					Texts: []string{"end"},
				},
			},
		},
	}
}
