package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/ignore"
)

func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".driftwatchignore")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeIgnoreFile: %v", err)
	}
	return path
}

func TestLoad_MissingFile(t *testing.T) {
	rules, err := ignore.Load("/nonexistent/.driftwatchignore")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if !rules.Empty() {
		t.Error("expected empty rules for missing file")
	}
}

func TestLoad_CommentsAndBlanks(t *testing.T) {
	path := writeIgnoreFile(t, "# this is a comment\n\n   \n")
	rules, err := ignore.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rules.Empty() {
		t.Error("expected empty rules")
	}
}

func TestMatch_ExactTypeAndPrefix(t *testing.T) {
	path := writeIgnoreFile(t, "aws_s3_bucket/my-bucket-")
	rules, _ := ignore.Load(path)

	if !rules.Match("aws_s3_bucket", "my-bucket-prod") {
		t.Error("expected match for exact type and prefix")
	}
	if rules.Match("aws_s3_bucket", "other-bucket") {
		t.Error("expected no match for different id")
	}
	if rules.Match("aws_instance", "my-bucket-prod") {
		t.Error("expected no match for different type")
	}
}

func TestMatch_WildcardType(t *testing.T) {
	path := writeIgnoreFile(t, "*/legacy-")
	rules, _ := ignore.Load(path)

	if !rules.Match("aws_instance", "legacy-host") {
		t.Error("expected match for wildcard type")
	}
	if !rules.Match("aws_s3_bucket", "legacy-data") {
		t.Error("expected match for wildcard type with different resource type")
	}
	if rules.Match("aws_instance", "prod-host") {
		t.Error("expected no match for non-matching prefix")
	}
}

func TestMatch_WildcardID(t *testing.T) {
	path := writeIgnoreFile(t, "aws_iam_role/*")
	rules, _ := ignore.Load(path)

	if !rules.Match("aws_iam_role", "any-role-id") {
		t.Error("expected match for wildcard id")
	}
	if rules.Match("aws_instance", "any-role-id") {
		t.Error("expected no match for different type with wildcard id")
	}
}

func TestMatch_MalformedLineSkipped(t *testing.T) {
	path := writeIgnoreFile(t, "no-slash-here\naws_s3_bucket/valid-")
	rules, _ := ignore.Load(path)

	if !rules.Match("aws_s3_bucket", "valid-bucket") {
		t.Error("expected valid pattern to still match")
	}
	if rules.Match("no-slash-here", "anything") {
		t.Error("malformed line should have been skipped")
	}
}
