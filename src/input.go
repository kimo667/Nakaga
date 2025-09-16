package main

import (
	"bufio"
	"fmt"
	"strings"
)

// clamp bornes une valeur entière entre min et max.
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// readLine lit une ligne avec un prompt, sans le \n final.
func readLine(r *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

// readChoice lit une ligne et renvoie en minuscules (pour comparer facilement).
func readChoice(r *bufio.Reader) string {
	return strings.ToLower(strings.TrimSpace(readLine(r, "> ")))
}
