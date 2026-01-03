# config

[![Go Reference][godoc-badge]][godoc-link]
[![Go Report Card][goreport-badge]][goreport-link]
[![License][license-badge]][license-link]

`config` provides **strict, crash-safe configuration file I/O** for Go projects.

It does **not** implement parsers or codecs. Instead, it defines a small, consistent
set of rules for loading and saving configuration files using existing format libraries.

## Install

```bash
go get github.com/BinGo-Lab-Team/config
```

## Packages

```
config/
├── jsonio/   # JSON configuration I/O
├── tomlio/   # TOML configuration I/O (strict, unknown keys rejected)
├── yamlio/   # YAML configuration I/O
├── xmlio/    # XML configuration I/O
```

Each package is format-specific and self-contained.

## API

All `*io` packages expose the same minimal API:

```go
func Load[T any](path string, cfg *T) error
func Save[T any](path string, cfg T) error
```

### Load

- Loads configuration from a file path
- Decodes exactly one document
- Returns errors on:
    - file not found
    - syntax errors
    - format-specific violations

### Save

- Writes configuration atomically:
    - temporary file
    - fsync
    - atomic rename
- Never leaves a partially written file
- Does not synchronize concurrent writes

## Usage

### TOML

```go
package main

import (
	"log"
	"errors"
	
	"github.com/BinGo-Lab-Team/config/tomlio"
)

type Config struct {
    Name string `toml:"name"`
    Port int    `toml:"port"`
}

func main() {
	var cfg Config
	if err := tomlio.Load("config.toml", &cfg); err != nil {
		if errors.Is(err, tomlio.ErrUnknownKeys) {
			// handle unknown keys
		}
		log.Fatal(err)
	}

	cfg.Port = 8081
	_ = tomlio.Save("config.toml", cfg)
}
```

## Notes

- File-level APIs only (no `io.Reader` / `io.Writer` wrappers)
- No schema, defaults, or hot-reload logic
- Intended for **application configuration**, not generic serialization

## License

This project released under [MIT License](./LICENSE)

[godoc-badge]: https://pkg.go.dev/badge/github.com/BinGo-Lab-Team/config.svg
[godoc-link]: https://pkg.go.dev/github.com/BinGo-Lab-Team/config

[goreport-badge]: https://goreportcard.com/badge/github.com/BinGo-Lab-Team/config
[goreport-link]: https://goreportcard.com/report/github.com/BinGo-Lab-Team/config

[license-badge]: https://img.shields.io/github/license/BinGo-Lab-Team/config
[license-link]: https://github.com/BinGo-Lab-Team/config/blob/main/LICENSE