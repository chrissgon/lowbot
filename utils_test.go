package lowbot

import (
	"testing"
)

func TestUtils_ParseTemplate(t *testing.T) {
	expect := "One\nTwo"
	have := ParseTemplate([]string{"One", "Two"})

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestUtils_IsURL(t *testing.T) {
	if !IsURL("http://www.google.com.br") {
		t.Errorf("is a valid url")
	}
	if IsURL("text") {
		t.Errorf("is a invalid url")
	}
}
