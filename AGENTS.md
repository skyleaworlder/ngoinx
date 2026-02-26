# AGENTS.md

## Cursor Cloud specific instructions

**ngoinx** is a lightweight Go reverse proxy / load balancer. Single binary, no Docker, no databases.

### Quick reference

| Action | Command |
|---|---|
| Build | `go build -o ngoinx ./ngoinx.go` |
| Run | `go run ngoinx.go` (reads `ngoinx.template.json`) |
| Lint | `go vet ./...` |
| Unit tests | `go test ./src/test/... -v` |

### Key caveats

- **Go module target is 1.15** (`go.mod`), but the code compiles fine with Go 1.22+. Do not upgrade the module version without need.
- **`Test_ConsistHash` has a pre-existing nil-pointer failure** — the test creates a `ConsistHash` without setting its logger. Run regexp/filepath tests separately with `-run Test_regexp` to avoid this.
- **Backend test servers** live in `src/test/server/test_*.py` (Flask). Install Flask (`pip3 install flask`) and start them before running ngoinx end-to-end. They are not needed for unit tests.
- **Path stripping**: ngoinx strips the `proxy.src` prefix before forwarding to backends. Backend `test_10081.py` registers route `/api/v1/test` but receives requests at `/` after stripping — this causes expected 404s. Backend `test_10082.py` registers `/` and works correctly. This is by design.
- **Config file**: `ngoinx.template.json` defines services on ports 10080 and 30080. Backends are expected on 10081-10086 and 30081.
- **Log and static dirs** (`log/`, `static/v1/`, `static/v3/`) must exist before running ngoinx — it writes log files and caches static resources there.
