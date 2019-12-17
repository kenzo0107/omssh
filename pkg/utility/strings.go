package utility

import "strings"

// ConvNewline : Convert LF code to specified string
func ConvNewline(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
	).Replace(str)
}
