// Command papercheck calculates the text similarity of two UTF-8 documents.
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run executes the command-line workflow. args must contain the original
// document, the suspected copy, and the result file, in that order.
func Run(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("用法: main.exe <原文绝对路径> <抄袭版绝对路径> <答案绝对路径>")
	}

	originalPath, candidatePath, resultPath := args[0], args[1], args[2]
	if err := EnsureOutputDistinct(resultPath, originalPath, candidatePath); err != nil {
		return err
	}

	original, err := ReadText(originalPath)
	if err != nil {
		return err
	}
	candidate, err := ReadText(candidatePath)
	if err != nil {
		return err
	}

	score := Similarity(Normalize(original), Normalize(candidate))
	if err := WriteResult(resultPath, score); err != nil {
		return err
	}
	return nil
}
