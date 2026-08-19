package pdf

import (
	"errors"
	"strings"
)

var (
	// ErrUnsupportedLanguage is returned when text is neither English nor Hindi.
	ErrUnsupportedLanguage = errors.New("unsupported or undetectable document language")
)

// DetectLanguage analyzes the text extracted from a PDF and determines if it is in English ("en") or Hindi ("hi").
func DetectLanguage(text string) (string, error) {
	trimmed := strings.TrimSpace(text)
	if len(trimmed) == 0 {
		return "", ErrUnsupportedLanguage
	}

	var devanagariCount, latinCount int

	for _, r := range trimmed {
		// Devanagari Unicode Block (Hindi): U+0900 to U+097F
		if r >= 0x0900 && r <= 0x097F {
			devanagariCount++
		} else if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			latinCount++
		}
	}

	totalRecognized := devanagariCount + latinCount
	if totalRecognized == 0 {
		return "", ErrUnsupportedLanguage
	}

	if devanagariCount > latinCount {
		return "hi", nil
	} else if latinCount > 0 {
		return "en", nil
	}

	return "", ErrUnsupportedLanguage
}
