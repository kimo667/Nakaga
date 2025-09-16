package main

import (
	"bufio"
	"fmt"
	"strings"
)

func readLine(r *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func readChoice(r *bufio.Reader) string {
	return strings.ToLower(strings.TrimSpace(readLine(r, "> ")))
}
