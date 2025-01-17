package utils

import "regexp"

func RemoveCotTag(s string) string {
	re := regexp.MustCompile(`<thinking>[\s\S]*?</thinking>`)
	return re.ReplaceAllString(s, "")
}
