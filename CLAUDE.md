# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository layout

Athena is a collection of independently-versioned Go libraries. **Every top-level directory is its own Go module** with its own `go.mod`/`go.sum`; there is no root `go.mod` and no `go.work`. Cross-module dependencies are resolved through published pseudo-versions in `go.sum`, not via `replace` directives or relative paths.

Modules and what they provide:

- `ast` — boolean-expression evaluator with Go-like syntax plus `${var}` references, `map[index]` membership, string slicing, and user-supplied functions. Has `STRICT` and `COMPATIBLE` modes (the latter accepts `=`, `&`, `|` in place of `==`, `&&`, `||`). The `ast/go/` subtree is a vendored fork of the Go 1.17 stdlib `go/ast`, `go/parser`, `go/token`, etc., modified to add `ast.Variable` and the extended operators — do not treat it as ordinary stdlib; changes there ripple through the parser.
- `bitmap` — fixed-size bitmap over `[]byte`.
- `broadcast` — fan-out of values to multiple goroutines via `sync.Cond`, with a fixed-capacity meta buffer that gates writes (data writes are dropped until `cap(meta)` metas are written).
- `consistent_hash` — consistent-hash ring built on a caller-supplied `ISortable` (no thread safety).
- `easybits` — struct (un)marshalling for bit-packed binary formats; depends on `easyerrors` and `easyio`. `Marshal` is a TODO.
- `easyerrors` — `HandleMultiError(f, errs...)` short-circuit helper.
- `easyhttpclient` — fluent HTTP client (`.M().T().H().P().B().Do()`); package name is `httpClient`, not the directory name.
- `easyio` — `EasyReader`/`EasyWriter` wrappers adding `ReadFull`, `ReadN`, `ReadAll`, `WriteFull`.
- `easypool` — `sync.Pool`-backed `Buffer` and `[]Buffer` pools.
- `easysyntax` — small syntax helpers (`DoLoop`, etc.); originally named differently (see "Module path quirks").
- `purl` — faster `net/url` replacement with optional whitelist/blacklist key indexing; depends on `easypool`.
- `ring_buffer` — ring-buffer cache built via a fluent builder (`NewRingBuffer(n).Array()|.List()|.Block().Build()`).

## Common commands

All commands run **inside a module directory**, since each module is independent:

```bash
cd <module>
go build ./...
go test ./...
go test -run TestName ./...      # single test
go test -race ./...
go vet ./...
go mod tidy
```

There is no Makefile, lint config, or CI workflow. To exercise every module from the repo root:

```bash
for d in ast bitmap broadcast consistent_hash easybits easyerrors easyhttpclient easyio easypool easysyntax purl ring_buffer; do (cd "$d" && go test ./...); done
```

## Module path quirks (read before editing imports or `go.mod`)

The repo is mid-migration from `github.com/SmartBrave/Athena/...` to `github.com/sbraveyoung/Athena/...`. Both paths are live in the tree:

- `SmartBrave/Athena/...`: `ast`, `bitmap`, `broadcast`, `consistent_hash`, `easybits`, `easyerrors`, `easyhttpclient`, `easyio`, `ring_buffer`
- `sbraveyoung/Athena/...`: `easypool`, `easysyntax`, `purl`

Imports inside a module must match its own `go.mod` path. When adding a cross-module import, check the consumer's `go.mod` first — for example, `easybits` (SmartBrave) imports `SmartBrave/Athena/easyerrors`, while `purl` (sbraveyoung) imports `sbraveyoung/Athena/easypool`. Don't "fix" one side without the other; bumping a dependency requires a tagged release of the dependency module first.

Also note that the package name in source does not always match the directory:

- `consistent_hash/` → `package consistentHash`
- `easyhttpclient/` → `package httpClient`
- `easypool/` → `package pool`

## Conventions

- Tests sit next to source as `*_test.go`; there is no separate test directory.
- Builders use short fluent method names (e.g. `easyhttpclient`'s `.M().T().H().P().B().Do()`, `ring_buffer`'s `.Array().Block().Build()`). Preserve this style when extending those APIs.
- The `ast` README (Chinese) is the canonical spec for the expression grammar — consult it before adding operators or changing precedence, and keep `STRICT` vs `COMPATIBLE` semantics in sync with the table there.
