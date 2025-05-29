package lowbot

import (
	"path/filepath"
	"reflect"
	"testing"
)

const (
	BUTTON = "button"
	TEXT   = "text"
	FILE   = "file"
)

var (
	BUTTONS   = []string{"button"}
	TO_MOCK   = NewWho("1", "chris")
	FROM_MOCK = NewWho("2", "amanda")
)

func TestInteraction_NewInteractionMessageButton(t *testing.T) {
	have := NewInteractionMessageButton(BUTTONS, TEXT)
	have.SetFrom(FROM_MOCK)
	have.SetTo(TO_MOCK)

	if TO_MOCK != have.To {
		t.Fatal(FormatTestError(TO_MOCK, have.To))
	}

	if FROM_MOCK != have.From {
		t.Fatal(FormatTestError(FROM_MOCK, have.From))
	}

	if MESSAGE_BUTTON != have.Type {
		t.Error(FormatTestError(MESSAGE_BUTTON, have.Type))
	}

	if TEXT != have.Parameters.Text {
		t.Error(FormatTestError(TEXT, have.Parameters.Text))
	}

	if len(BUTTONS) != len(have.Parameters.Buttons) {
		t.Error(FormatTestError(BUTTONS, have.Parameters.Buttons))
	}
}

func TestInteraction_NewInteractionMessageFile(t *testing.T) {
	have := NewInteractionMessageFile(FILE, "", TEXT)
	have.SetFrom(FROM_MOCK)
	have.SetTo(TO_MOCK)

	if TO_MOCK != have.To {
		t.Fatal(FormatTestError(TO_MOCK, have.To))
	}

	if FROM_MOCK != have.From {
		t.Fatal(FormatTestError(FROM_MOCK, have.From))
	}

	if MESSAGE_FILE != have.Type {
		t.Error(FormatTestError(MESSAGE_FILE, have.Type))
	}

	if TEXT != have.Parameters.Text {
		t.Error(FormatTestError(TEXT, have.Parameters.Text))
	}

	abs, _ := filepath.Abs(FILE)

	if abs != have.Parameters.File.GetFile().Path {
		t.Error(FormatTestError(abs, have.Parameters.File.GetFile().Path))
	}
}

func TestInteraction_NewInteractionMessageText(t *testing.T) {
	have := NewInteractionMessageText(TEXT)
	have.SetFrom(FROM_MOCK)
	have.SetTo(TO_MOCK)

	if TO_MOCK != have.To {
		t.Fatal(FormatTestError(TO_MOCK, have.To))
	}

	if FROM_MOCK != have.From {
		t.Fatal(FormatTestError(FROM_MOCK, have.From))
	}

	if MESSAGE_TEXT != have.Type {
		t.Error(FormatTestError(MESSAGE_TEXT, have.Type))
	}

	if TEXT != have.Parameters.Text {
		t.Error(FormatTestError(TEXT, have.Parameters.Text))
	}
}

func TestInteraction_SetReplier(t *testing.T) {
	interaction := NewInteractionMessageText(TEXT)

	interaction.SetReplier(WHO_MOCK)

	have := interaction.Replier
	expect := WHO_MOCK

	if !reflect.DeepEqual(expect, have) {
		t.Error(FormatTestError(expect, have))
	}
}

func TestInteraction_SetFrom(t *testing.T) {
	interaction := NewInteractionMessageText(TEXT)

	interaction.SetFrom(WHO_MOCK)

	have := interaction.From
	expect := WHO_MOCK

	if !reflect.DeepEqual(expect, have) {
		t.Error(FormatTestError(expect, have))
	}
}

func TestInteraction_SetTo(t *testing.T) {
	interaction := NewInteractionMessageText(TEXT)

	interaction.SetTo(WHO_MOCK)

	have := interaction.To
	expect := WHO_MOCK

	if !reflect.DeepEqual(expect, have) {
		t.Error(FormatTestError(expect, have))
	}
}
