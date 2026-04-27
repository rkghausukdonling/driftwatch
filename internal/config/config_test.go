package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/driftwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	tmp := filepath.Join(t.TempDir(), "driftwatch.yaml")
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return tmp
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Provider != "aws" {
		t.Errorf("expected default provider 'aws', got %q", cfg.Provider)
	}
	if cfg.OutputFmt != "table" {
		t.Errorf("expected default output_format 'table', got %q", cfg.OutputFmt)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := `
provider: gcp
region: us-central1
statefile: ./terraform.tfstate
output_format: json
ignore:
  - aws_s3_bucket.logs
tags:
  env: staging
`
	path := writeTemp(t, raw)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Provider != "gcp" {
		t.Errorf("expected provider 'gcp', got %q", cfg.Provider)
	}
	if cfg.OutputFmt != "json" {
		t.Errorf("expected output_format 'json', got %q", cfg.OutputFmt)
	}
	if len(cfg.Ignore) != 1 || cfg.Ignore[0] != "aws_s3_bucket.logs" {
		t.Errorf("unexpected ignore list: %v", cfg.Ignore)
	}
}

func TestLoad_InvalidProvider(t *testing.T) {
	raw := `provider: digitalocean\noutput_format: table\n`
	path := writeTemp(t, raw)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for unsupported provider, got nil")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/driftwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestValidate_InvalidOutputFormat(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.OutputFmt = "xml"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid output_format, got nil")
	}
}
