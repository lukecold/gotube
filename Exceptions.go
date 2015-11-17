package main

import (
	"fmt"
	. "strings"
)

type AgeRestrictedError struct {}

type MissingFieldsError struct {
	_fields []string
}

type NoMatchingVideoError struct {
	_quality   string
	_extension string
}

type PatternNotFoundError struct {
	_pattern string
}

type UnmatchedBracketsError struct{}

func (AgeRestrictedError) Error() string {
	return fmt.Sprintf("Error: this page is age-restricted")
}

func (e MissingFieldsError) Error() string {
	return fmt.Sprintf("Error: detected missing field(s): %v", Join(e._fields, ", "))
}

func (e NoMatchingVideoError) Error() string {
	errorMessage := "Error: no matching videos with"
	if e._quality != "" {
		errorMessage += " quality=" + e._quality
	}
	if e._extension != "" {
		errorMessage += " extension=" + e._extension
	}
	return fmt.Sprintf(errorMessage)
}

func (e PatternNotFoundError) Error() string {
	return fmt.Sprintf("Error: pattern \"%v\" not found", e._pattern)
}

func (UnmatchedBracketsError) Error() string {
	return fmt.Sprintf("Error: unmatched brackets")
}