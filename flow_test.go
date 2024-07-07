package lowbot

import (
	"reflect"
	"testing"
)

func TestFlow_NewFlow(t *testing.T) {
	expect := newFlowMock()
	have, err := NewFlow("./mocks/flow.yaml")

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, *have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFlow_StartFlow(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()

	expect := flow.Steps["init"]
	have := flow.Current

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFlow_NextFlow(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, ""))

	expect := flow.Steps["audio"]
	have := flow.Current

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestFlow_NoHasNext(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Current.Next = nil

	if !flow.NoHasNext() {
		t.Errorf(FormatTestError(true, false))
	}
}

func TestFlow_AddResponse(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, ""))

	expect := 1
	have := len(flow.Current.Responses)

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

}
func TestFlow_AddResponseValue(t *testing.T) {
	interaction := NewInteractionMessageText(CHANNEL_MOCK, DESTINATION_MOCK, SENDER_MOCK, "Response")
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(interaction)

	expect := interaction
	have := flow.Current.Responses[0]

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
