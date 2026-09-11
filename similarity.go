package main

// CountBigrams returns the frequency of every adjacent two-rune feature.
// Keeping counts makes repeated passages contribute proportionally.
func CountBigrams(text []rune) map[[2]rune]int {
	counts := make(map[[2]rune]int)
	for index := 0; index+1 < len(text); index++ {
		feature := [2]rune{text[index], text[index+1]}
		counts[feature]++
	}
	return counts
}

// CountCharacters returns the frequency of each rune.
func CountCharacters(text []rune) map[rune]int {
	counts := make(map[rune]int)
	for _, char := range text {
		counts[char]++
	}
	return counts
}

// Similarity calculates a multiset Dice similarity in [0, 1]. Bigrams are
// used for normal documents; character frequencies provide a stable fallback
// when either document has fewer than two normalized characters.
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
