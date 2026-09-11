package main

import "unicode"

// Normalize keeps letters and digits, removes whitespace and punctuation,
// and lowercases Unicode letters. Runes are used so Chinese text is handled
// as characters rather than as individual UTF-8 bytes.
func Normalize(text string) []rune {
	// len(text) is a byte count, so this may reserve extra capacity for UTF-8
	// input, but it avoids a second full rune conversion and remains linear.
	result := make([]rune, 0, len(text))
	for _, char := range text {
		if char == '\uFEFF' {
			continue
		}
		if unicode.IsLetter(char) {
			result = append(result, unicode.ToLower(char))
			continue
		}
		if unicode.IsDigit(char) {
			result = append(result, char)
		}
	}
	return result
}
