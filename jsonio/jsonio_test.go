package jsonio_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/BinGo-Lab-Team/config/jsonio"
)

type testConfig struct {
	Name string `json:"name"`
	Port int    `json:"port"`
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	orig := testConfig{
		Name: "example",
		Port: 8080,
	}

	// Save config
	if err := jsonio.Save(path, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load config
	var loaded testConfig
	if err := jsonio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != orig {
		t.Fatalf("loaded config mismatch: got %+v, want %+v", loaded, orig)
	}
}

func TestLoadFileNotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "not-exist.json")

	var cfg testConfig
	err := jsonio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}
}

func TestSaveOverwriteIsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	first := testConfig{
		Name: "first",
		Port: 1,
	}
	second := testConfig{
		Name: "second",
		Port: 2,
	}

	// Initial save
	if err := jsonio.Save(path, first); err != nil {
		t.Fatalf("initial Save failed: %v", err)
	}

	// Overwrite
	if err := jsonio.Save(path, second); err != nil {
		t.Fatalf("overwrite Save failed: %v", err)
	}

	var loaded testConfig
	if err := jsonio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != second {
		t.Fatalf("overwrite failed: got %+v, want %+v", loaded, second)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// Write invalid JSON manually
	if err := os.WriteFile(path, []byte(`{ invalid json }`), 0o644); err != nil {
		t.Fatalf("write invalid json: %v", err)
	}

	var cfg testConfig
	err := jsonio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}
