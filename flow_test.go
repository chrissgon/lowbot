package lowbot

import (
	"reflect"
	"testing"
)

func TestNewFlow(t *testing.T) {
	expect := NewMock()
	have, err := NewFlow("./mocks/flow.yaml")

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, *have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestStartFlow(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()

	expect := flow.Steps["init"]
	have := flow.Current

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestNextFlow(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(NewInteractionMessageText(CHANNELID, SESSIONID, ""))

	expect := flow.Steps["audio"]
	have := flow.Current

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestNoHasNext(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Current.Next = nil

	if !flow.NoHasNext() {
		t.Errorf(FormatTestError(true, false))
	}
}

func TestAddResponse(t *testing.T) {
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(NewInteractionMessageText(CHANNELID, SESSIONID, ""))

	expect := 1
	have := len(flow.Current.Responses)

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}

}
func TestAddResponseValue(t *testing.T) {
	in := NewInteractionMessageText(CHANNELID, SESSIONID, "Response")
	flow, _ := NewFlow("./mocks/flow.yaml")
	flow.Start()
	flow.Next(in)

	expect := in
	have := flow.Current.Responses[0]

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(FormatTestError(expect, have))
	}
}

func NewMock() Flow {
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
				Action: "Audio",
				Parameters: StepParameters{
					Audio: "./mocks/music.mp3",
				},
			},
			"button": {
				Next: map[string]string{
					"default": "document",
				},
				Action: "Button",
				Parameters: StepParameters{
					Buttons: []string{"yes", "no"},
					Texts: []string{"buttons here"},
				},
			},
			"document": {
				Next: map[string]string{
					"default": "image",
				},
				Action: "Document",
				Parameters: StepParameters{
					Document: "./mocks/features.txt",
				},
			},
			"image": {
				Next: map[string]string{
					"default": "text",
				},
				Action: "Image",
				Parameters: StepParameters{
					Image: "./mocks/image.jpg",
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
				Action: "Video",
				Parameters: StepParameters{
					Video: "./mocks/video.mp4",
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
