package sanitizer

import "github.com/microcosm-cc/bluemonday"

var sanitizer *bluemonday.Policy

func init() {
	sanitizer = bluemonday.UGCPolicy()
}

func Sanitize(s string) string {
	return sanitizer.Sanitize(s)
}
