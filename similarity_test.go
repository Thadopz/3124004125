package main

import (
	"math"
	"testing"
)

func TestCountBigramsKeepsRepeatedFeatures(t *testing.T) {
	got := CountBigrams([]rune("aaaa"))
	if got[[2]rune{'a', 'a'}] != 3 {
		t.Fatalf("CountBigrams() = %#v, want aa count 3", got)
	}
}

func TestSimilarityIdenticalText(t *testing.T) {
	if got := Similarity([]rune("今天是星期天"), []rune("今天是星期天")); got != 1 {
		t.Fatalf("Similarity() = %v, want 1", got)
	}
}

func TestSimilarityBothEmpty(t *testing.T) {
	if got := Similarity(nil, nil); got != 1 {
		t.Fatalf("Similarity() = %v, want 1", got)
	}
}

func TestSimilarityOneEmpty(t *testing.T) {
	if got := Similarity(nil, []rune("文本")); got != 0 {
		t.Fatalf("Similarity() = %v, want 0", got)
	}
}

func TestSimilarityNoCommonBigrams(t *testing.T) {
	if got := Similarity([]rune("甲乙丙丁"), []rune("戊己庚辛")); got != 0 {
		t.Fatalf("Similarity() = %v, want 0", got)
	}
}

func TestSimilarityHandCalculatedPointSixSeven(t *testing.T) {
	got := Similarity([]rune("甲乙丙丁"), []rune("甲乙丙戊"))
	if math.Abs(got-2.0/3.0) > 1e-12 {
		t.Fatalf("Similarity() = %v, want 0.666...", got)
	}
}

func TestSimilarityInsertion(t *testing.T) {
	got := Similarity([]rune("今天是星期天"), []rune("今天是星期天天气晴"))
	if got <= 0 || got >= 1 {
		t.Fatalf("Similarity() = %v, want between 0 and 1", got)
	}
}

func TestSimilarityDeletion(t *testing.T) {
	got := Similarity([]rune("今天是星期天天气晴"), []rune("今天是星期天"))
	if got <= 0 || got >= 1 {
		t.Fatalf("Similarity() = %v, want between 0 and 1", got)
	}
}

func TestSimilarityReplacement(t *testing.T) {
	got := Similarity([]rune("今天是星期天"), []rune("今天是工作日"))
	if got <= 0 || got >= 1 {
		t.Fatalf("Similarity() = %v, want between 0 and 1", got)
	}
}

func TestSimilarityReorderedTextRetainsLocalFeatures(t *testing.T) {
	got := Similarity([]rune("甲乙丙丁戊己"), []rune("丁戊己甲乙丙"))
	if got <= 0 || got >= 1 {
		t.Fatalf("Similarity() = %v, want between 0 and 1", got)
	}
}

func TestSimilarityRepeatedPassageUsesMultiplicity(t *testing.T) {
	got := Similarity([]rune("甲乙甲乙"), []rune("甲乙甲乙甲乙"))
	if math.Abs(got-0.75) > 1e-12 {
		t.Fatalf("Similarity() = %v, want 0.75", got)
	}
}

func TestSimilaritySingleCharacterFallback(t *testing.T) {
	got := Similarity([]rune("甲"), []rune("甲乙"))
	if math.Abs(got-2.0/3.0) > 1e-12 {
		t.Fatalf("Similarity() = %v, want 0.666...", got)
	}
}

func TestSimilaritySingleDifferentCharacters(t *testing.T) {
	if got := Similarity([]rune("甲"), []rune("乙")); got != 0 {
		t.Fatalf("Similarity() = %v, want 0", got)
	}
}

func TestSimilarityNormalizesBeforeComparison(t *testing.T) {
	if got := Similarity(Normalize("今天是星期天。"), Normalize("今天是星期天")); got != 1 {
		t.Fatalf("Similarity() = %v, want 1", got)
	}
}

func TestSimilarityUsesCharacterFallbackWhenCandidateIsShort(t *testing.T) {
	got := Similarity([]rune("甲乙丙"), []rune("甲"))
	if math.Abs(got-0.5) > 1e-12 {
		t.Fatalf("Similarity() = %v, want 0.5", got)
	}
}
