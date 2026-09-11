package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// ReadText reads one UTF-8 document. A UTF-8 BOM at the beginning is ignored.
func ReadText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件 %q 失败: %w", path, err)
	}
	if !utf8.Valid(data) {
		return "", fmt.Errorf("读取文件 %q 失败: 文件不是有效的 UTF-8 文本", path)
	}
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	return string(data), nil
}

// WriteResult writes a score in the format required by the assignment.
func WriteResult(path string, score float64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("写入答案文件 %q 失败: %w", path, err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "%.2f\n", score); err != nil {
		return fmt.Errorf("写入答案文件 %q 失败: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("保存答案文件 %q 失败: %w", path, err)
	}
	return nil
}

// EnsureOutputDistinct prevents the result file from overwriting either input.
// It checks normalized paths as well as filesystem identity, so aliases such
// as "dir\\..\\orig.txt" and symlinks are handled when they resolve locally.
func EnsureOutputDistinct(resultPath string, inputPaths ...string) error {
	resultNormalized, err := normalizedPath(resultPath)
	if err != nil {
		return fmt.Errorf("检查答案路径失败: %w", err)
	}

	resultInfo, resultStatErr := os.Stat(resultPath)
	for _, inputPath := range inputPaths {
		inputNormalized, err := normalizedPath(inputPath)
		if err != nil {
			return fmt.Errorf("检查输入路径失败: %w", err)
		}
		if pathsEqual(resultNormalized, inputNormalized) {
			return fmt.Errorf("答案文件不能覆盖输入文件: %q", inputPath)
		}

		inputInfo, inputStatErr := os.Stat(inputPath)
		if resultStatErr == nil && inputStatErr == nil && os.SameFile(resultInfo, inputInfo) {
			return fmt.Errorf("答案文件不能覆盖输入文件: %q", inputPath)
		}
	}
	return nil
}

func normalizedPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	absPath = filepath.Clean(absPath)
	if resolved, err := filepath.EvalSymlinks(absPath); err == nil {
		absPath = filepath.Clean(resolved)
	}
	return absPath, nil
}

func pathsEqual(first, second string) bool {
	if filepath.Separator == '\\' {
		return strings.EqualFold(first, second)
	}
	return first == second
}
