package tomlio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// ErrUnknownKeys indicates that the TOML file contains
// keys not defined in the target struct.
var ErrUnknownKeys = errors.New("unknown keys")

// formatKeys formats TOML keys into a human-readable string.
// Used for error reporting only.
func formatKeys(keys []toml.Key) string {
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, strings.Join(k, "."))
	}
	return strings.Join(parts, ", ")
}

// Load loads a TOML configuration file from path into cfg.
//
// cfg must be a pointer.
// Unknown keys in the file result in ErrUnknownKeys.
func Load[T any](path string, cfg *T) error {
	path = filepath.Clean(path)

	meta, err := toml.DecodeFile(path, cfg)
	if err != nil {
		// err may wrap fs.ErrNotExist if the file does not exist
		return fmt.Errorf("decode %s: %w", path, err)
	}

	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		return fmt.Errorf("%w: %s", ErrUnknownKeys, formatKeys(undecoded))
	}

	return nil
}

// Save writes cfg as TOML to path.
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
	tmp, err := os.CreateTemp(dir, "*.toml.tmp")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpPath := tmp.Name()

	// Encode TOML into the temporary file.
	enc := toml.NewEncoder(tmp)
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
