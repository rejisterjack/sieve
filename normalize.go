package sieve

import (
	"strings"
	"unicode"
)

// zero-width and related format characters often used to obfuscate injections.
var zeroWidth = map[rune]struct{}{
	'\u200b': {}, // zero width space
	'\u200c': {}, // zero width non-joiner
	'\u200d': {}, // zero width joiner
	'\ufeff': {}, // BOM / zero width no-break space
	'\u2060': {}, // word joiner
	'\u180e': {}, // mongolian vowel separator
}

// normalize prepares input for pattern matching: strip zero-width chars,
// collapse whitespace runs, and lowercase.
func normalize(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	prevSpace := false
	for _, r := range input {
		if _, ok := zeroWidth[r]; ok {
			continue
		}
		if unicode.IsSpace(r) {
			if prevSpace {
				continue
			}
			b.WriteByte(' ')
			prevSpace = true
			continue
		}
		prevSpace = false
		b.WriteRune(unicode.ToLower(r))
	}
	return strings.TrimSpace(b.String())
}
