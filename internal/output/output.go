// Package output provides writer selection for scan results.
package output

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Format represents a supported output format.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ErrUnknownFormat is returned when an unsupported format is requested.
type ErrUnknownFormat struct {
	Name string
}

func (e *ErrUnknownFormat) Error() string {
	return fmt.Sprintf("unknown output format %q: supported formats are text, json", e.Name)
}

// ParseFormat validates and normalises a format string.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(strings.TrimSpace(s))) {
	case FormatText:
		return FormatText, nil
	case FormatJSON:
		return FormatJSON, nil
	default:
		return "", &ErrUnknownFormat{Name: s}
	}
}

// WriterFor returns an io.WriteCloser for the given path.
// If path is "-" or empty, os.Stdout is returned (Close is a no-op).
func WriterFor(path string) (io.WriteCloser, error) {
	if path == "" || path == "-" {
		return nopCloser{os.Stdout}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output: create file %q: %w", path, err)
	}
	return f, nil
}

// nopCloser wraps a writer with a no-op Close method.
type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }
