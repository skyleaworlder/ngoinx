# AGENTS.md — ngoinx

## Project Overview

**ngoinx** is a lightweight nginx-like reverse proxy written in Go. It was built as a learning project to explore the core mechanics behind nginx.

- **Module:** `github.com/skyleaworlder/ngoinx`
- **License:** MIT
- **Language:** Go (go.mod declares 1.15; any modern Go toolchain 1.15+ works)

### Core Features

1. **Load Balancing** — Automatic algorithm selection based on target count:
   - 1–3 targets → Weighted Round-Robin
   - 4+ targets → Consistent Hashing
2. **Reverse Proxy** — HTTP request forwarding based on JSON-configured source-to-target mappings.
3. **Static/Dynamic Separation** — Static resources (`.html`, `.css`, `.js`, `.jpg`, `.png`, `.ico`) are cached locally; dynamic requests are proxied to upstream targets.
4. **JSON Configuration** — All settings live in `ngoinx.template.json`, parsed via `gjson`.

## Architecture & Project Structure

```
ngoinx.go                        # Entry point — reads config, inits load balancers, starts server
ngoinx.template.json             # Configuration file (Service → Proxy → Target hierarchy)

src/
  config/
    GLOVAR.go                    # Global variable: Svc []Service
    Service.go                   # Service struct + Unmarshal(gjson.Result)
    Proxy.go                     # Proxy and Target structs + Unmarshal(gjson.Result)

  ldbls/                         # Load balancers
    Ldbls.go                     # LoadBalancer interface, LdblserMap (global), LdblserMapStuffer
    ConsistHash.go               # Consistent hashing (sha1-based, uses go-sortedmap)
    WeightedRoundRobin.go        # Weighted round-robin

  server/
    Server.go                    # Server struct, NewNgoinxServer, Serve, handlerGenerator

  staticfs/                      # Static file system / caching
    Folder.go                    # Folder struct (implements http.FileSystem), Open/Create/clean
    FolderManager.go             # FolderManager struct, Init/Serve/StripURLPathPrefix/clean

  utils/
    Config.go                    # ReadConfig — reads JSON config, populates config.Svc
    Logger.go                    # Loggerable interface, LoggerConfig, LoggerGenerator
    Hash.go                      # GetHashID — sha1 hash to uint64
    DeltaUint64.go               # Ring distance calculation for consistent hashing
    Marshal.go                   # Marshaler / Unmarshaler interfaces
    RegexUtils.go                # StaticSuffixDetermine, StripURLPathPrefix, GetStaticFileName

  test/                          # Tests (non-standard location — not alongside source files)
    consist_hash_test.go         # Consistent hashing tests
    ldbls_test.go                # Load balancer map stuffer + server integration test
    filepath_test.go             # Static folder / filepath tests
    regexp_test.go               # Regex utility tests
    marshal_test.go              # Config unmarshal tests

doc/
  config-format.md               # Configuration file format documentation
  roadmap.md                     # Feature roadmap

log/                             # Log output directory (generated at runtime)
static/                          # Static file cache directory (generated at runtime)
```

### Key Interfaces

| Interface | Package | Purpose |
|-----------|---------|---------|
| `LoadBalancer` | `src/ldbls` | Contract for load balancing algorithms: `Init`, `GetAddr`, `SetLogger` |
| `Loggerable` | `src/utils` | Any struct needing configurable logging implements `SetLogger(*LoggerConfig)` |
| `Unmarshaler` | `src/utils` | Config structs implement `Unmarshal(gjson.Result)` for JSON deserialization |

### Global State

- `config.Svc []config.Service` — populated by `utils.ReadConfig`, holds all service definitions.
- `ldbls.LdblserMap map[string]LoadBalancer` — populated by `ldbls.LdblserMapStuffer`, maps source paths to load balancers.

## Dependencies

| Dependency | Purpose |
|------------|---------|
| `github.com/sirupsen/logrus` | Structured logging with per-component log files |
| `github.com/tidwall/gjson` | Fast JSON parsing for config files |
| `github.com/umpc/go-sortedmap` | Sorted map used by ConsistHash for ordered key iteration |

Install dependencies:

```bash
go mod download
```

## Development Environment

- **Go toolchain:** 1.15+ (tested with 1.22)
- **No external services** required — ngoinx is self-contained.
- **EditorConfig** is provided (`.editorconfig`):
  - JSON files: 2-space indent, UTF-8, LF line endings
  - Go files: 4-space indent

## How to Build & Run

```bash
# Run directly
go run ngoinx.go

# Or build and run
go build -o ngoinx . && ./ngoinx
```

The server reads `ngoinx.template.json` from the current directory. It will:
1. Parse the config and populate `config.Svc`.
2. Initialize load balancers for each proxy route (`ldbls.LdblserMapStuffer`).
3. Start an HTTP server on each configured `listen` port.
4. Start the static folder manager's periodic cleanup goroutine.

Log files are written to `./log/` (one per service, plus per-component logs for load balancers and static folders).

### Configuration

Edit `ngoinx.template.json` to configure services. Structure:

```json
{
  "service": [
    {
      "listen": 10080,
      "log": "./log/",
      "static": "./static/v1/",
      "proxy": [
        {
          "src": "/api/v1/test",
          "target": [
            { "dst": "http://127.0.0.1:10081", "weight": 3 },
            { "dst": "http://127.0.0.1:10082", "weight": 2 }
          ]
        }
      ]
    }
  ]
}
```

See `doc/config-format.md` for the full specification.

## How to Test

Tests are located in `src/test/` (non-standard Go convention — all test files are in a single directory).

```bash
# Run all tests (from the test directory)
cd src/test && go test -v ./...

# Individual test suites
cd src/test && go test -v -run Test_ConsistHash        # Consistent hashing
cd src/test && go test -v -run Test_regexp              # Regex utilities
cd src/test && go test -v -run Test_Marshal_Unmarshal   # Config parsing
cd src/test && go test -v -run Test_filepath            # Static folder paths
```

### Important Notes on Tests

- **Relative paths:** Some tests reference config files via relative paths (`../../ngoinx.template.json` or `../config/ngoinx.template.json`). Always run tests from the `src/test/` directory.
- **Integration test:** `Test_LdblsMapStuffer` in `ldbls_test.go` starts the full server and sleeps for 1000 seconds. It is a long-running integration test, not a quick unit test. Avoid running it in CI without a timeout.
- **Test package:** All test files declare `package main`, not their respective source packages.

## Code Conventions

### Naming

- **Go files:** PascalCase (e.g., `ConsistHash.go`, `FolderManager.go`, `GLOVAR.go`)
- **Packages:** lowercase, short names (`ldbls`, `staticfs`, `config`, `utils`, `server`)
- **Exported types/functions:** PascalCase per Go convention
- **Constructors:** `NewDefault*` for default constructors, `New*` for parameterized constructors

### Patterns

- **Interface-driven design:** Core behaviors are defined as interfaces (`LoadBalancer`, `Loggerable`, `Unmarshaler`). New implementations should satisfy the relevant interface.
- **Per-component logging:** Each major struct (load balancers, static folders, server) gets its own log file via `SetLogger(*LoggerConfig)`.
- **Config deserialization:** Config structs implement `Unmarshal(gjson.Result)` to populate themselves from parsed JSON.
- **Error handling:** Log warnings/fatals via logrus, return `error` up the call chain. No panics.
- **Global state:** `config.Svc` and `ldbls.LdblserMap` are module-level globals populated at startup.

### Adding a New Load Balancer

1. Create a new file in `src/ldbls/` (e.g., `RandomSelect.go`).
2. Define a struct that implements the `LoadBalancer` interface (`Init`, `GetAddr`, `SetLogger`).
3. Add a `NewDefault*` constructor.
4. Update the selection logic in `LdblserMapStuffer` in `src/ldbls/Ldbls.go`.

## Key Design Decisions

- **Automatic load balancer selection:** The algorithm is chosen based on the number of targets per proxy route (1–3 → round-robin, 4+ → consistent hashing). This logic lives in `ldbls.LdblserMapStuffer`.
- **Static resource detection:** File extensions are matched via regex (`(.html|.css|.js|.jpg|.png|.ico)$`). Defined in `utils.StaticSuffixDetermine`.
- **Cache cleanup:** The `FolderManager` runs a goroutine that periodically (every 60s) deletes empty files and files older than 1 minute from static cache folders.
- **Reverse proxy path stripping:** The source prefix (`service[i].proxy[j].src`) is stripped from the request path before forwarding to the target. E.g., `/api/v1/food/1.js` with src `/api/v1/food` forwards as `/1.js`.
