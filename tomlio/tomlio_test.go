package tomlio_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/BinGo-Lab-Team/config/tomlio"
	"github.com/BurntSushi/toml"
)

type testConfig struct {
	Name string `toml:"name"`
	Port int    `toml:"port"`
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	orig := testConfig{
		Name: "example",
		Port: 8080,
	}

	if err := tomlio.Save(path, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loaded testConfig
	if err := tomlio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != orig {
		t.Fatalf("loaded config mismatch: got %+v, want %+v", loaded, orig)
	}
}

func TestLoadFileNotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "not-exist.toml")

	var cfg testConfig
	err := tomlio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}
}

func TestLoadUnknownKeys(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	// Write TOML with an unknown key.
	content := `
name = "example"
port = 8080
extra = "unexpected"
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write toml: %v", err)
	}

	var cfg testConfig
	err := tomlio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected ErrUnknownKeys, got nil")
	}

	if !errors.Is(err, tomlio.ErrUnknownKeys) {
		t.Fatalf("expected ErrUnknownKeys, got %v", err)
	}
}

func TestSaveOverwriteIsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	first := testConfig{
		Name: "first",
		Port: 1,
	}
	second := testConfig{
		Name: "second",
		Port: 2,
	}

	if err := tomlio.Save(path, first); err != nil {
		t.Fatalf("initial Save failed: %v", err)
	}

	if err := tomlio.Save(path, second); err != nil {
		t.Fatalf("overwrite Save failed: %v", err)
	}

	var loaded testConfig
	if err := tomlio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != second {
		t.Fatalf("overwrite failed: got %+v, want %+v", loaded, second)
	}
}

func TestLoadInvalidTOML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	// Invalid TOML syntax.
	if err := os.WriteFile(path, []byte(`= invalid toml =`), 0o644); err != nil {
		t.Fatalf("write invalid toml: %v", err)
	}

	var cfg testConfig
	err := tomlio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}

	// Ensure the error comes from TOML decoding.
	var parseErr toml.ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected toml parse error, got %v", err)
	}
}
