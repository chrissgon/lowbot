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
	matched, _ := regexp.MatchString(`^(?:https?://)?(?:[^/.\s]+\.).*`, str)
	return matched
}

func PrintLog(msg string) {
	if DEBUG {
		log.Print(msg)
	}
}
