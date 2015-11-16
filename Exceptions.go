package main

import (
	. "strings"
	"fmt"
)

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

type MissingFieldsError struct {
	fields []string
}

func (e MissingFieldsError) Error() string {
	return fmt.Sprintf("Error: detected missing field(s): %v", Join(e.fields, ", "))
}