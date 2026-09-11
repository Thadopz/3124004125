package main

// CountBigrams 统计每个相邻双 rune 特征的出现次数。
// 保留出现次数，使重复片段按其频率参与相似度计算。
func CountBigrams(text []rune) map[[2]rune]int {
	counts := make(map[[2]rune]int)
	for index := 0; index+1 < len(text); index++ {
		feature := [2]rune{text[index], text[index+1]}
		counts[feature]++
	}
	return counts
}

// CountCharacters 统计每个 rune 的出现次数。
func CountCharacters(text []rune) map[rune]int {
	counts := make(map[rune]int)
	for _, char := range text {
		counts[char]++
	}
	return counts
}

// Similarity 计算 [0, 1] 范围内的多重集 Dice 相似度。正常长度的文档使用二元组；
// 任一文档规范化后少于两个字符时，稳定地回退到字符频次计算。
func Similarity(original, candidate []rune) float64 {
	if sameRunes(original, candidate) {
		return 1
	}
	if len(original) == 0 || len(candidate) == 0 {
		return 0
	}
	if len(original) < 2 || len(candidate) < 2 {
		return diceCharacters(original, candidate)
	}
	return diceBigrams(CountBigrams(original), CountBigrams(candidate), len(original)-1, len(candidate)-1)
}

func sameRunes(first, second []rune) bool {
	if len(first) != len(second) {
		return false
	}
	for index := range first {
		if first[index] != second[index] {
			return false
		}
	}
	return true
}

func diceBigrams(first, second map[[2]rune]int, firstTotal, secondTotal int) float64 {
	intersection := 0
	for feature, count := range first {
		if otherCount := second[feature]; otherCount < count {
			intersection += otherCount
		} else {
			intersection += count
		}
	}
	return diceScore(intersection, firstTotal, secondTotal)
}

func diceCharacters(first, second []rune) float64 {
	firstCounts := CountCharacters(first)
	secondCounts := CountCharacters(second)
	intersection := 0
	for char, count := range firstCounts {
		if otherCount := secondCounts[char]; otherCount < count {
			intersection += otherCount
		} else {
			intersection += count
		}
	}
	return diceScore(intersection, len(first), len(second))
}

func diceScore(intersection, firstTotal, secondTotal int) float64 {
	denominator := firstTotal + secondTotal
	if denominator == 0 {
		return 1
	}
	return float64(2*intersection) / float64(denominator)
}
