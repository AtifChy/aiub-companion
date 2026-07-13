package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	hookSrc := filepath.Join("scripts", "pre-commit.sh")
	hookDst := filepath.Join(".git", "hooks", "pre-commit")

	// Ensure .git/hooks directory exists
	err := os.MkdirAll(filepath.Dir(hookDst), 0o755)
	if err != nil {
		fmt.Printf("Error creating hooks directory: %v\n", err)
		os.Exit(1)
	}

	srcFile, err := os.Open(hookSrc)
	if err != nil {
		fmt.Printf("Error opening source hook: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = srcFile.Close() }()

	dstFile, err := os.OpenFile(hookDst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		fmt.Printf("Error creating destination hook: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Printf("Error copying hook: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Git pre-commit hook installed successfully!")
}
