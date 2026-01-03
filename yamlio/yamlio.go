package yamlio

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load loads a YAML configuration file from path into cfg.
//
// cfg must be a pointer.
// The function performs a single YAML decode from the file.
func Load[T any](path string, cfg *T) error {
	path = filepath.Clean(path)

	fp, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer func() { _ = fp.Close() }()

	if err := yaml.NewDecoder(fp).Decode(cfg); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}

	return nil
}

// Save writes cfg as YAML to path.
//
// Save is crash-safe:
//   - data is written to a temporary file
//   - the file is synced
//   - the file is atomically renamed to the target path
//
// The parent directory is created if it does not exist.
// Concurrent writes to the same path are not synchronized.
func Save[T any](path string, cfg T) error {
	path = filepath.Clean(path)
	dir := filepath.Dir(path)

	// Ensure parent directory exists.
	if err := os.MkdirAll(dir, 0o0755); err != nil {
		return fmt.Errorf("make dir: %w", err)
	}

	// Create temporary file in the same directory.
	tmp, err := os.CreateTemp(dir, "*.yml.tmp")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpPath := tmp.Name()

	enc := yaml.NewEncoder(tmp)
	defer func() { _ = enc.Close() }()

	// Encode YAML into the temporary file.
	if err := enc.Encode(cfg); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write %s: %w", tmpPath, err)
	}

	// Flush file contents to disk.
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("sync %s: %w", tmpPath, err)
	}

	// Close before renaming.
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close %s: %w", tmpPath, err)
	}

	// Atomically replace the target file.
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename %s -> %s: %w", tmpPath, path, err)
	}

	// Best-effort directory sync.
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		_ = d.Close()
	}

	return nil
}
