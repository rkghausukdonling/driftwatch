package output_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/driftwatch/internal/output"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  output.Format
	}{
		{"text", output.FormatText},
		{"TEXT", output.FormatText},
		{"json", output.FormatJSON},
		{"JSON", output.FormatJSON},
		{" json ", output.FormatJSON},
	}
	for _, tc := range cases {
		got, err := output.ParseFormat(tc.input)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
	var uf *output.ErrUnknownFormat
	if ok := isErrUnknownFormat(err, &uf); !ok {
		t.Fatalf("expected *ErrUnknownFormat, got %T", err)
	}
	if uf.Name != "xml" {
		t.Errorf("ErrUnknownFormat.Name = %q; want %q", uf.Name, "xml")
	}
}

func isErrUnknownFormat(err error, target **output.ErrUnknownFormat) bool {
	if e, ok := err.(*output.ErrUnknownFormat); ok {
		*target = e
		return true
	}
	return false
}

func TestWriterFor_Stdout(t *testing.T) {
	w, err := output.WriterFor("-")
	if err != nil {
		t.Fatalf("WriterFor("-") error: %v", err)
	}
	defer w.Close()
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriterFor_Empty(t *testing.T) {
	w, err := output.WriterFor("")
	if err != nil {
		t.Fatalf("WriterFor("") error: %v", err)
	}
	defer w.Close()
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriterFor_File(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	w, err := output.WriterFor(path)
	if err != nil {
		t.Fatalf("WriterFor(%q) error: %v", path, err)
	}
	w.Close()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist at %q: %v", path, err)
	}
}

func TestWriterFor_InvalidPath(t *testing.T) {
	_, err := output.WriterFor("/nonexistent/dir/out.txt")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
