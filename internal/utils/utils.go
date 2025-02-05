package utils

import (
	"regexp"
	"strings"
)

func GetCotContext(s string) string {
	re := regexp.MustCompile(`<thinking>[\s\S]*?</thinking>`)
	s = re.FindString(s)
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(s, "</thinking>"), "<thinking>"))
}
func RemoveCotTag(s string) string {
	re := regexp.MustCompile(`<thinking>[\s\S]*?</thinking>`)
	return strings.TrimSpace(re.ReplaceAllString(s, ""))
}
