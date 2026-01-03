package yamlio_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	// adjust import path to your module
	"github.com/BinGo-Lab-Team/config/yamlio"
)

type testConfig struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	orig := testConfig{
		Name: "example",
		Port: 8080,
	}

	if err := yamlio.Save(path, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loaded testConfig
	if err := yamlio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != orig {
		t.Fatalf("loaded config mismatch: got %+v, want %+v", loaded, orig)
	}
}

func TestLoadFileNotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "not-exist.yaml")

	var cfg testConfig
	err := yamlio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}
}

func TestSaveOverwriteIsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	first := testConfig{
		Name: "first",
		Port: 1,
	}
	second := testConfig{
		Name: "second",
		Port: 2,
	}

	if err := yamlio.Save(path, first); err != nil {
		t.Fatalf("initial Save failed: %v", err)
	}

	if err := yamlio.Save(path, second); err != nil {
		t.Fatalf("overwrite Save failed: %v", err)
	}

	var loaded testConfig
	if err := yamlio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != second {
		t.Fatalf("overwrite failed: got %+v, want %+v", loaded, second)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	// Invalid YAML content.
	if err := os.WriteFile(path, []byte(`name: [unclosed`), 0o644); err != nil {
		t.Fatalf("write invalid yaml: %v", err)
	}

	var cfg testConfig
	err := yamlio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}
