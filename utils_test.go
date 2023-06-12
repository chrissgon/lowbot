package lowbot

import (
	"testing"
)

func TestParseTemplate(t *testing.T) {
	expect := "One\nTwo"
	have := ParseTemplate([]string{"One", "Two"})

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
func TestInt64ToString(t *testing.T) {
	expect := "47"
	have := Int64ToString(47)

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestStringToInt64(t *testing.T) {
	expect := int64(47)
	have := StringToInt64("47")

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestIsURL(t *testing.T) {
	if !IsURL("http://www.google.com.br") {
		t.Errorf("is a valid url")
	}
	if IsURL("text") {
		t.Errorf("is a invalid url")
	}
}
