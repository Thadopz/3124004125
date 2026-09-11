package main

import "testing"

func BenchmarkSimilarityLargeText(b *testing.B) {
	text := Normalize("今天是星期天，天气晴，今天晚上我要去看电影。论文查重性能测试。")
	large := make([]rune, 0, len(text)*10000)
	for index := 0; index < 10000; index++ {
		large = append(large, text...)
	}
	// Make the documents similar but different so the benchmark measures
	// feature counting and intersection instead of the equality fast path.
	candidate := append(append([]rune(nil), large...), '改')
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		Similarity(large, candidate)
	}
}
