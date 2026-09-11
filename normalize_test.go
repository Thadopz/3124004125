package main

import (
	"reflect"
	"testing"
)

func TestNormalizeRemovesWhitespaceAndPunctuation(t *testing.T) {
	got := Normalize("今天是星期天，天气晴！")
	want := []rune("今天是星期天天气晴")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Normalize() = %q, want %q", string(got), string(want))
	}
}

func TestNormalizeLowercasesUnicodeLetters(t *testing.T) {
	got := Normalize("Go语言 ABC １２３")
	want := []rune("go语言abc１２３")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Normalize() = %q, want %q", string(got), string(want))
	}
}

func TestNormalizeSkipsBOM(t *testing.T) {
	got := Normalize("\uFEFF原文")
	want := []rune("原文")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Normalize() = %q, want %q", string(got), string(want))
	}
}

func TestNormalizeKeepsCombiningLetterBaseOnly(t *testing.T) {
	got := Normalize("e\u0301")
	want := []rune("e")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Normalize() = %q, want %q", string(got), string(want))
	}
}
