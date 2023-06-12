package lowbot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FormatTestError(expect, have interface{}) string {
	return fmt.Sprintf("\nexpect: %v, \nhave: %v", expect, have)
}

func ParseTemplate(texts []string) string {
	return strings.Join(texts, "\n")
}

func Int64ToString(number int64) string {
	return fmt.Sprintf("%d", number)
}

func StringToInt64(str string) int64 {
	number, _ := strconv.ParseInt(str, 10, 64)
	return number
}

func IsURL(str string) bool {
	matched, _ := regexp.MatchString(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`, str)
	return matched
}

func NewError(fn string, err error) error {
	return fmt.Errorf("%s: %w", fn, err)
}
