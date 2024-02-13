package store

import "strings"

const (
	ErrTokenTypeNotString   = "token of type string required"
	DefaultUnredactedLength = 4
	DefaultRedactedToken    = "*"
)

// InterfaceIsString checks if an interface underlying type is a string
func InterfaceIsString(a any) (bool, string) {
	switch v := a.(type) {
	case string:
		return true, v
	default:
		return false, ""
	}
}

// Redact censors sensitive token before printing them in logs or as a response.
func Redact(s string) string {
	redacted := ""
	oLength := len(s)
	redacted += s[:DefaultUnredactedLength]

	redacted += strings.Repeat(DefaultRedactedToken, oLength-len(redacted))
	return redacted
}
