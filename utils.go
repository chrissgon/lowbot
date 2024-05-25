package lowbot

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func FormatTestError(expect, have any) string {
	return fmt.Sprintf("\nexpect: %v, \nhave: %v", expect, have)
}

func ParseTemplate(texts []string) string {
	return strings.Join(texts, "\n")
}

func IsURL(str string) bool {
	matched, _ := regexp.MatchString(`^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$`, str)
	return matched
}

func printLog(msg string) {
	if DEBUG {
		log.Print(msg)
	}
}
