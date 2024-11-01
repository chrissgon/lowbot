package lowbot

const (
	BUTTON = "button"
	TEXT   = "text"
	FILE   = "file"
)

var (
	BUTTONS          = []string{"button"}
	DESTINATION_MOCK = NewWho("1", "chris")
	SENDER_MOCK      = NewWho("2", "amanda")
)

// func TestInteraction_NewInteractionMessageButton(t *testing.T) {
// 	have := NewInteractionMessageButton(DESTINATION_MOCK, SENDER_MOCK, BUTTONS, TEXT)

// 	if DESTINATION_MOCK != have.Destination {
// 		t.Error(FormatTestError(DESTINATION_MOCK, have.Destination))
// 	}

// 	if SENDER_MOCK != have.Sender {
// 		t.Error(FormatTestError(SENDER_MOCK, have.Sender))
// 	}

// 	if MESSAGE_BUTTON != have.Type {
// 		t.Error(FormatTestError(MESSAGE_BUTTON, have.Type))
// 	}

// 	if TEXT != have.Parameters.Text {
// 		t.Error(FormatTestError(TEXT, have.Parameters.Text))
// 	}

// 	if len(BUTTONS) != len(have.Parameters.Buttons) {
// 		t.Error(FormatTestError(BUTTONS, have.Parameters.Buttons))
// 	}
// }

// func TestInteraction_NewInteractionMessageFile(t *testing.T) {
// 	have := NewInteractionMessageFile(DESTINATION_MOCK, SENDER_MOCK, FILE, TEXT)

// 	if DESTINATION_MOCK != have.Destination {
// 		t.Error(FormatTestError(DESTINATION_MOCK, have.Destination))
// 	}

// 	if SENDER_MOCK != have.Sender {
// 		t.Error(FormatTestError(SENDER_MOCK, have.Sender))
// 	}

// 	if MESSAGE_FILE != have.Type {
// 		t.Error(FormatTestError(MESSAGE_FILE, have.Type))
// 	}

// 	if TEXT != have.Parameters.Text {
// 		t.Error(FormatTestError(TEXT, have.Parameters.Text))
// 	}

// 	abs, _ := filepath.Abs(FILE)

// 	if abs != have.Parameters.File.GetFile().Path {
// 		t.Error(FormatTestError(abs, have.Parameters.File.GetFile().Path))
// 	}
// }

// func TestInteraction_NewInteractionMessageText(t *testing.T) {
// 	have := NewInteractionMessageText(DESTINATION_MOCK, SENDER_MOCK, TEXT)

// 	if DESTINATION_MOCK != have.Destination {
// 		t.Error(FormatTestError(DESTINATION_MOCK, have.Destination))
// 	}

// 	if SENDER_MOCK != have.Sender {
// 		t.Error(FormatTestError(SENDER_MOCK, have.Sender))
// 	}

// 	if MESSAGE_TEXT != have.Type {
// 		t.Error(FormatTestError(MESSAGE_TEXT, have.Type))
// 	}

// 	if TEXT != have.Parameters.Text {
// 		t.Error(FormatTestError(TEXT, have.Parameters.Text))
// 	}
// }

// func TestInteraction_SetReplier(t *testing.T) {
// 	interaction := NewInteractionMessageText(DESTINATION_MOCK, SENDER_MOCK, TEXT)

// 	interaction.SetReplier(WHO_MOCK)

// 	have := interaction.Replier
// 	expect := WHO_MOCK

// 	if !reflect.DeepEqual(expect, have) {
// 		t.Error(FormatTestError(expect, have))
// 	}
// }

// func TestInteraction_SetDestination(t *testing.T) {
// 	interaction := NewInteractionMessageText(DESTINATION_MOCK, SENDER_MOCK, TEXT)

// 	interaction.SetDestination(WHO_MOCK)

// 	have := interaction.Destination
// 	expect := WHO_MOCK

// 	if !reflect.DeepEqual(expect, have) {
// 		t.Error(FormatTestError(expect, have))
// 	}
// }
