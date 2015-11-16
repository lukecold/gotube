package main

import "fmt"

type PatternNotFoundError struct {
	pattern string
}

func (e PatternNotFoundError) Error() string {
	return fmt.Sprintf("Error: pattern \"%v\" not found", e.pattern)
}

type UnmatchedBracketsError struct{}

func (UnmatchedBracketsError) Error() string {
	return fmt.Sprintf("Error: unmatched brackets")
}
