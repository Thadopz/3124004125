package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadTextRemovesBOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "原文.txt")
	if err := os.WriteFile(path, append([]byte{0xEF, 0xBB, 0xBF}, []byte("内容")...), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := ReadText(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "内容" {
		t.Fatalf("ReadText() = %q, want 内容", got)
	}
}

func TestReadTextRejectsInvalidUTF8(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.txt")
	if err := os.WriteFile(path, []byte{0xff, 0xfe}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadText(path); err == nil || !strings.Contains(err.Error(), "UTF-8") {
		t.Fatalf("ReadText() error = %v, want UTF-8 error", err)
	}
}

func TestWriteResultUsesTwoDecimalPlaces(t *testing.T) {
	path := filepath.Join(t.TempDir(), "答案.txt")
	if err := WriteResult(path, 2.0/3.0); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "0.67\n" {
		t.Fatalf("result = %q, want 0.67\\n", data)
	}
}

func TestEnsureOutputDistinctRejectsSamePath(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "原文.txt")
	if err := os.WriteFile(input, []byte("内容"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureOutputDistinct(input, input); err == nil {
		t.Fatal("EnsureOutputDistinct() accepted identical input and output")
	}
}

func TestEnsureOutputDistinctAllowsNewOutput(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "原文.txt")
	output := filepath.Join(directory, "答案.txt")
	if err := os.WriteFile(input, []byte("内容"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureOutputDistinct(output, input); err != nil {
		t.Fatalf("EnsureOutputDistinct() = %v, want nil", err)
	}
}

func TestEnsureOutputDistinctRejectsHardLink(t *testing.T) {
	directory := t.TempDir()
	input := filepath.Join(directory, "原文.txt")
	alias := filepath.Join(directory, "答案.txt")
	if err := os.WriteFile(input, []byte("内容"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(input, alias); err != nil {
		t.Skipf("当前文件系统不支持硬链接: %v", err)
	}
	if err := EnsureOutputDistinct(alias, input); err == nil {
		t.Fatal("EnsureOutputDistinct() accepted a hard-link alias")
	}
}

func TestRunWritesExpectedResult(t *testing.T) {
	directory := t.TempDir()
	originalPath := filepath.Join(directory, "原文.txt")
	candidatePath := filepath.Join(directory, "抄袭版.txt")
	resultPath := filepath.Join(directory, "答案.txt")
	if err := os.WriteFile(originalPath, []byte("甲乙丙丁"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(candidatePath, []byte("甲乙丙戊"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{originalPath, candidatePath, resultPath}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "0.67\n" {
		t.Fatalf("result = %q, want 0.67\\n", data)
	}
}

func TestRunRejectsWrongArgumentCount(t *testing.T) {
	if err := Run([]string{"only-one"}); err == nil {
		t.Fatal("Run() accepted an invalid argument count")
	}
}

func TestRunRejectsOutputAliasingOriginal(t *testing.T) {
	directory := t.TempDir()
	originalPath := filepath.Join(directory, "orig.txt")
	candidatePath := filepath.Join(directory, "copy.txt")
	if err := os.WriteFile(originalPath, []byte("甲乙"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(candidatePath, []byte("甲乙"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{originalPath, candidatePath, filepath.Join(directory, ".", "orig.txt")}); err == nil {
		t.Fatal("Run() accepted an aliased output path")
	}
}

func TestRunMissingInputKeepsExistingResult(t *testing.T) {
	directory := t.TempDir()
	originalPath := filepath.Join(directory, "orig.txt")
	resultPath := filepath.Join(directory, "answer.txt")
	if err := os.WriteFile(originalPath, []byte("甲乙"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(resultPath, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Run([]string{originalPath, filepath.Join(directory, "missing.txt"), resultPath})
	if err == nil {
		t.Fatal("Run() accepted a missing input")
	}
	data, readErr := os.ReadFile(resultPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != "keep\n" {
		t.Fatalf("existing result changed to %q", data)
	}
}

func TestRunInvalidUTF8KeepsExistingResult(t *testing.T) {
	directory := t.TempDir()
	originalPath := filepath.Join(directory, "orig.txt")
	candidatePath := filepath.Join(directory, "candidate.txt")
	resultPath := filepath.Join(directory, "answer.txt")
	if err := os.WriteFile(originalPath, []byte("甲乙"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(candidatePath, []byte{0xff, 0xfe}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(resultPath, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Run([]string{originalPath, candidatePath, resultPath})
	if err == nil {
		t.Fatal("Run() accepted invalid UTF-8")
	}
	data, readErr := os.ReadFile(resultPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != "keep\n" {
		t.Fatalf("existing result changed to %q", data)
	}
}

func TestRunRejectsDirectoryAsResult(t *testing.T) {
	directory := t.TempDir()
	originalPath := filepath.Join(directory, "orig.txt")
	candidatePath := filepath.Join(directory, "candidate.txt")
	resultPath := filepath.Join(directory, "result-directory")
	if err := os.WriteFile(originalPath, []byte("甲乙"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(candidatePath, []byte("甲乙"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(resultPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{originalPath, candidatePath, resultPath}); err == nil {
		t.Fatal("Run() accepted a directory as result path")
	}
}
