package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b strings.Builder

	var pending rune
	hasPending := false
	escaped := false

	for _, r := range s {
		if escaped {
			if !(isDigit(r) || r == '\\') {
				return "", ErrInvalidString
			}
			if hasPending {
				b.WriteRune(pending)
			}
			pending, hasPending = r, true
			escaped = false
			continue
		}

		if r == '\\' {
			if hasPending {
				b.WriteRune(pending)
				hasPending = false
			}
			escaped = true
			continue
		}

		if isDigit(r) {
			if !hasPending {
				return "", ErrInvalidString
			}
			n := int(r - '0')
			for range n {
				b.WriteRune(pending)
			}
			hasPending = false
			continue
		}

		if hasPending {
			b.WriteRune(pending)
		}
		pending, hasPending = r, true
	}

	if escaped {
		return "", ErrInvalidString
	}
	if hasPending {
		b.WriteRune(pending)
	}

	return b.String(), nil
}

func isDigit(r rune) bool { return r >= '0' && r <= '9' }
