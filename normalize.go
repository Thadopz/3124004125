package main

import "unicode"

// Normalize 保留字母和数字，移除空白与标点，并将 Unicode 字母转换为小写。
// 使用 rune 逐个处理字符，使中文按字符而不是按 UTF-8 字节处理。
func Normalize(text string) []rune {
	// len(text) 表示字节数，因此对 UTF-8 输入可能会预留多余容量；这样可以避免再次完整
	// 转换为 rune，并保持线性复杂度。
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
