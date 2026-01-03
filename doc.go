// Package config provides strict, crash-safe configuration file I/O.
//
// This package does not implement parsers or codecs.
// Instead, it defines a consistent, file-level contract for loading and
// saving configuration files using established format libraries.
//
// Subpackages under config/*io provide format-specific implementations
// (JSON, TOML, YAML, XML), all exposing the same minimal API:
//
//	Load[T any](path string, cfg *T) error
//	Save[T any](path string, cfg T) error
//
// Design principles:
//
//   - Configuration files are treated as critical inputs.
//     Errors must be explicit and observable.
//   - Writes are crash-safe and never leave partially written files.
//   - Behavior is format-explicit; no automatic format detection is performed.
//   - The API operates at the file level; io.Reader/io.Writer are intentionally
//     not abstracted.
//
// This package is intended for application configuration,
// not for general-purpose data serialization.
package config
