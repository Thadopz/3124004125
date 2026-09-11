package main

import "testing"

func BenchmarkSimilarityLargeText(b *testing.B) {
	text := Normalize("今天是星期天，天气晴，今天晚上我要去看电影。论文查重性能测试。")
	large := make([]rune, 0, len(text)*10000)
	for index := 0; index < 10000; index++ {
		large = append(large, text...)
	}
	// 构造内容相似但不完全相同的文档，使基准测试实际测量特征统计和交集计算，避免命中相等文本的快速路径。
	candidate := append(append([]rune(nil), large...), '改')
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		Similarity(large, candidate)
	}
}
