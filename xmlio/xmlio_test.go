package xmlio_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/BinGo-Lab-Team/config/xmlio"
)

type testConfig struct {
	XMLName struct{} `xml:"config"`
	Name    string   `xml:"name"`
	Port    int      `xml:"port"`
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.xml")

	orig := testConfig{
		Name: "example",
		Port: 8080,
	}

	if err := xmlio.Save(path, orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loaded testConfig
	if err := xmlio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != orig {
		t.Fatalf("loaded config mismatch: got %+v, want %+v", loaded, orig)
	}
}

func TestLoadFileNotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "not-exist.xml")

	var cfg testConfig
	err := xmlio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}
}

func TestSaveOverwriteIsAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.xml")

	first := testConfig{
		Name: "first",
		Port: 1,
	}
	second := testConfig{
		Name: "second",
		Port: 2,
	}

	if err := xmlio.Save(path, first); err != nil {
		t.Fatalf("initial Save failed: %v", err)
	}

	if err := xmlio.Save(path, second); err != nil {
		t.Fatalf("overwrite Save failed: %v", err)
	}

	var loaded testConfig
	if err := xmlio.Load(path, &loaded); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded != second {
		t.Fatalf("overwrite failed: got %+v, want %+v", loaded, second)
	}
}

func TestLoadInvalidXML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.xml")

	// Invalid XML content.
	if err := os.WriteFile(path, []byte(`<config><name></config>`), 0o644); err != nil {
		t.Fatalf("write invalid xml: %v", err)
	}

	var cfg testConfig
	err := xmlio.Load(path, &cfg)
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
}
