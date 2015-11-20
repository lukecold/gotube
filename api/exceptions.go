package gotube

import (
	"fmt"
	. "strings"
)

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

func (e MissingFieldsError) Error() string {
	return fmt.Sprintf("detected missing field(s): %v", Join(e._fields, ", "))
}

func (e NoMatchingVideoError) Error() string {
	errorMessage := "no matching videos with"
	if e._quality != "" {
		errorMessage += " quality=" + e._quality
	}
	if e._extension != "" {
		errorMessage += " extension=" + e._extension
	}
	return fmt.Sprintf(errorMessage)
}

func (e PatternNotFoundError) Error() string {
	return fmt.Sprintf("pattern \"%v\" not found", e._pattern)
}